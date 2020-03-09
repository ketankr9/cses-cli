package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/briandowns/spinner"
)

/*
api's
// mandatory
cses login
cses list
cses show 1068
cses solve 1068
cses submit 1068.filename.cpp

// TODO
cses result 1068

// optional
cses stat 1068
*/

type Session struct {
	Csrf   string       `json:"csrf"`
	User   string       `json:"username"`
	Cookie string       `json:"cookie"`
	Root   string       `json:"root"`
	Editor string       `json:"editor"`
	Lang string `json:"lang"`
	Github githubConfig `json:"github"`
}

// var cpptemplate = `
// #include<bits/stdc++.h>
// using namespace std;

// int main(){

// 	return 0;
// }
// `

func initSess(sess *Session) bool {
	os.MkdirAll(sess.Root, os.ModePerm)

	out, ok := cacheGet("config.json", sess.Root)
	if !ok {
		return false
	}

	json.Unmarshal(out, sess)
	if sess.Csrf == "" || sess.User == "" || sess.Cookie == "" {
		return false
	}

	return true
}

func login(sess *Session, pass string) bool {

	params := "csrf_token=" + sess.Csrf + "&nick=" + sess.User + "&pass=" + pass
	if loginRequest(params, sess.Cookie) == sess.User {
		return true
	}
	return false
}

func promtLogin(sess *Session) bool {
	scanner := bufio.NewScanner(os.Stdin)

	updateIfNew(scanner, &sess.User, "Username")

	fmt.Print("Password: ")
	scanner.Scan()
	PASSWORD := scanner.Text()

	sess.Cookie, sess.Csrf = newCookieCsrf()
	ok := login(sess, PASSWORD)

	if !ok {
		return false
	}

	updateConfig(sess)

	return true
}

func list(sess *Session) {

	doc, err := goquery.NewDocumentFromReader(listRequest(sess.Cookie))
	check(err)

	doc.Find(".task").Each(func(i int, s *goquery.Selection) {

		solved := "✘"

		a := s.Find("a")
		link, _ := a.Attr("href")
		taskNumber := link[17:]
		title := a.Text()

		hitRatio := strings.Split(s.Find("span").Text(), "/")
		n, err := strconv.ParseFloat(strings.TrimSpace(hitRatio[0]), 64)
		check(err)
		d, err := strconv.ParseFloat(strings.TrimSpace(hitRatio[1]), 64)
		check(err)

		percent := n * 100 / d

		s.Find(".task-score").Each(func(o int, k *goquery.Selection) {
			st, _ := k.Attr("class")
			if strings.Contains(st, "full") {
				solved = "✔"
			} else if "task-score icon " == st {
				solved = "-"
			}
		})

		fmt.Printf("\t%s [%s] %-25s (%.1f %%)\n", solved, taskNumber, title, percent)
	})

}

func printResult(link string, sess *Session) bool {

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = "PENDING "
	s.Start()

	for true {
		status, text, verdict := printResultRequest(link, sess.Cookie)
		s.Prefix = status + " "

		if status == "READY" || status == "COMPILE ERROR" || status == "" {
			s.Stop()
			fmt.Print(text)
			if verdict == "ACCEPTED" {
				return true
			}
			break
		}
	}
	return false
}

func submit(sourceFile string, sess *Session) {

	ext := filepath.Ext(sourceFile)

	opts := map[string]string{
		"csrf_token": sess.Csrf,
		"task":       strings.Split(filepath.Base(sourceFile), ".")[0],
		"lang":       extLangMap[ext],
		"type":       "course",
		"target":     "problemset",
		"option":     extOptionMap[ext],
	}

	link := submitRequest(opts, sourceFile, sess.Cookie)

	if verdict := printResult(link, sess); verdict && validGithubConfig(&sess.Github) {
		s := spinner.New(spinner.CharSets[36], 100*time.Millisecond)
		s.Prefix = "Committing to Github"
		s.Start()
		if ok := updateFile(sourceFile, &sess.Github); ok {
			s.Stop()
			fmt.Println("Github: "+sess.Github.SourceRepo+" ✔")
		} else {
			s.Stop()
			fmt.Println("Github: ✘")
		}
	}
}

func getTask(task string, sess *Session) (string, bool) {
	filename := task + ".task.html"
	path := filepath.Join(sess.Root, filename)

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Prefix = "Downloading "
		s.Start()
		defer s.Stop()

		text := downloadTask(task)
		if text == "" {
			return "", false
		}
		cacheSet(filename, text, sess.Root)
	}

	return getTaskFromCache(task, sess.Root), true
}

func show(task string, sess *Session) {
	text, exist := getTask(task, sess)
	if exist {
		fmt.Println(text)
	} else {
		fmt.Println("Task Doesn't Exist")
	}
}

func solve(task string, sess *Session) {

	text, exist := getTask(task, sess)
	if !exist {
		fmt.Println("Task Doesn't Exist")
	}

	filename := task + ".task" + langExtMap[sess.Lang]
	template := getTemplate(langExtMap[sess.Lang])

	writeCodeFile(filename, text, template)

	if sess.Editor == "" {
		scanner := bufio.NewScanner(os.Stdin)

		updateIfNew(scanner, &sess.Editor, "Editor")

		if sess.Editor == "" {
			fmt.Println("Editor still not configured")
			return
		}
		updateConfig(sess)
	}
	exec.Command(sess.Editor, filename).Output()
}

func configureGithub(sess *Session) {
	scanner := bufio.NewScanner(os.Stdin)

	updateIfNew(scanner, &sess.Github.Token, "Token")
	updateIfNew(scanner, &sess.Github.SourceRepo, "Repository")
	updateIfNew(scanner, &sess.Github.AuthorName, "Github Username")
	updateIfNew(scanner, &sess.Github.AuthorEmail, "Github Email")

	updateConfig(sess)
}

func showHelp(){
	fmt.Println("Usage:")

	fmt.Println("\tcses-cli login")
	fmt.Println("\tcses-cli list")
	fmt.Println("\tcses-cli show 1068")
	fmt.Println("\tcses-cli solve 1068")
	fmt.Println("\tcses-cli submit 1068.task.cpp")
	
	fmt.Println("Optional:")
	fmt.Println("\tcses-cli github")
}

// func stat(task string) {
// 	fmt.Println("#Todo")
// }

func main() {

	flag.Parse()

	if flag.NArg() == 0 {
		showHelp()
		return
	}

	sess := &Session{
		Lang: "cpp",
		Root: filepath.Join(UserHomeDir(), ".cses"),
	}

	isLogged := initSess(sess)

	switch flag.Arg(0) {
	case "login":
		if !promtLogin(sess) {
			fmt.Println("Login failed")
		} else {
			fmt.Println("Logged in successfully")
		}
	case "list":
		list(sess)
	case "show":
		if flag.NArg() < 2 {
			os.Exit(1)
		}
		show(flag.Arg(1), sess)
	case "solve":
		if flag.NArg() < 2 {
			os.Exit(1)
		}
		show(flag.Arg(1), sess)
		solve(flag.Arg(1), sess)
	case "submit":
		if flag.NArg() < 2 {
			os.Exit(1)
		}
		if !isLogged {
			fmt.Println("\tPlease login first :(")
			return
		}
		submit(flag.Arg(1), sess)
	case "github":
		configureGithub(sess)

	// case "stat":
	// 	fmt.Println("sw stat 1068")
	// 	stat("stat 1068")
	case "help":
	default:
		showHelp()
	}
}
