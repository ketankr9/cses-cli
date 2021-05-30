package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Problem struct {
	Task     string  `json:"task"`
	Title    string  `json:"title"`
	Solved   string  `json:"solved"`
	HitRatio float64 `json:"hitratio"`
}

func UpdateProblemListCache(cookie string) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = "Updating Problem List "
	s.Start()
	problems := ListRequest(cookie)
	cacheSet("problemList.json", problems, rootPath)
	s.Stop()
}

func ListProblems(cookie string) {
	var problems []*Problem
	cacheGet("problemList.json", &problems, rootPath)

	if len(problems) == 0 {
		UpdateProblemListCache(cookie)
		cacheGet("problemList.json", &problems, rootPath)
	}
	for _, v := range problems {
		fmt.Printf("\t%s [%s] %-25s (%.1f %%)\n", v.Solved, v.Task, v.Title, v.HitRatio)
	}
}

func ShowProblem(taskNumber string) {
	_, text, exist := getTask(taskNumber)
	if exist {
		fmt.Println(text)
	} else {
		fmt.Println("Task Doesn't Exist")
	}
}

func SolveProblem(taskNumber string, lang string, editor string) {

	title, text, exist := getTask(taskNumber)
	if !exist {
		fmt.Println("Task Doesn't Exist")
	}

	filename := taskNumber + "." + title + langExtMap[lang]
	template := getTemplate(langExtMap[lang])

	writeCodeFile(filename, text, template)

	_, err := exec.Command(editor, filename).Output()
	check(err)
}

func printResult(link string, cookie string) bool {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = "PENDING "
	s.Start()

	for true {
		status, text, verdict := printResultRequest(link, cookie)
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

func SubmitSolution(sourceFile string, cookie string, csrf string, ghConfig *GithubConfig) {

	ext := filepath.Ext(sourceFile)

	opts := map[string]string{
		"csrf_token": csrf,
		"task":       strings.Split(filepath.Base(sourceFile), ".")[0],
		"lang":       extLangMap[ext],
		"type":       "course",
		"target":     "problemset",
		"option":     extOptionMap[ext],
	}

	link := submitRequest(opts, sourceFile, cookie)

	if verdict := printResult(link, cookie); verdict {
		go updateListByTask(opts["task"], "✔", rootPath)
		if validGithubConfig(ghConfig) {
			s := spinner.New(spinner.CharSets[36], 100*time.Millisecond)
			s.Prefix = "Committing to Github"
			s.Start()
			if ok := updateFile(sourceFile, ghConfig); ok {
				s.Stop()
				fmt.Println("Github: " + ghConfig.SourceRepo + " ✔")
			} else {
				s.Stop()
				fmt.Println("Github: ✘")
			}
		}
	} else {
		updateListByTask(opts["task"], "✘", rootPath)
	}
}

func getTask(task string) (string, string, bool) {
	filename := task + ".task.html"
	path := filepath.Join(rootPath, filename)

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Prefix = "Downloading "

		s.Start()
		text := DownloadTask(task)
		s.Stop()

		if text == "" {
			return "", "", false
		}
		writeTask(filename, text, rootPath)
	}
	title, text := getTaskFromCache(task, rootPath)
	return title, text, true
}