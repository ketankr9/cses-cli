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

cses result 1068

// optional
cses stat 1068
*/
type Session struct {
	Csrf   string `json:"csrfToken"`
	User   string `json:"username"`
	Cookie string `json:"cookie"`
	Root   string `json:"root"`
	Editor string `json:"editor"`
}

var cpptemplate = `
#include<bits/stdc++.h>
using namespace std;

int main(){

	return 0;
}
`

func initSess(sess *Session) bool {
	sess.Root = filepath.Join(UserHomeDir(), ".cses")
	os.MkdirAll(sess.Root, os.ModePerm)

	out, ok := cacheGet("login.json", sess.Root)
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

func generateNewSess(sess *Session) {

	sess.Cookie, sess.Csrf = newCookieCsrf()

}

func promtLogin(sess *Session) bool {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Username: ")
	scanner.Scan()
	sess.User = scanner.Text()

	fmt.Print("Password: ")
	scanner.Scan()
	PASSWORD := scanner.Text()

	generateNewSess(sess)
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

func printResult(link string, sess *Session) {

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = "PENDING "
	s.Start()

	for true {
		status, text := printResultRequest(link, sess.Cookie)
		s.Prefix = status + " "

		if status == "READY" || status == "" {
			fmt.Println("\n" + text)
			break
		}
	}

	s.Stop()

}

func submit(filename string, sess *Session) {

	task, lang, option, exist := fileMeta(filename)
	if !exist {
		fmt.Println("File doesn't exist")
		return
	}

	opts := map[string]string{
		"csrf_token": sess.Csrf,
		"task":       task,
		"lang":       lang,
		"type":       "course",
		"target":     "problemset",
		"option":     option,
	}

	link := submitRequest(opts, filename, sess.Cookie)
	printResult(link, sess)
}

func getTask(task string, sess *Session) (string, bool) {
	filename := task + ".task.html"
	path := filepath.Join(sess.Root, filename)

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Prefix = "Downloading "
		s.Start()

		text := downloadTask(task)
		cacheSet(filename, text, sess.Root)

		s.Stop()
	}

	return getTaskFromCache(task, sess.Root), true
}

func show(task string, sess *Session) {
	text, ok := getTask(task, sess)
	if ok {
		fmt.Println(text)
	} else {
		fmt.Println("Task Doesn't Exist")
	}
}

func solve(task string, sess *Session) {

	text, _ := getTask(task, sess)
	filename := task + ".task.cpp"

	writeCodeFile(filename, text, cpptemplate)

	if sess.Editor == "" {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Editor: ")
		scanner.Scan()
		sess.Editor = scanner.Text()

		updateConfig(sess)
	}
	exec.Command(sess.Editor, filename).Output()
}

func stat(task string) {
	fmt.Println("#Todo")
}

func main() {

	flag.Parse()

	if flag.NArg() == 0 {
		os.Exit(1)
	}

	sess := &Session{}
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

	case "stat":
		fmt.Println("sw stat 1068")
		stat("stat 1068")

	default:
		fmt.Println("sw default")
	}
}
