package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Problem struct {
	Task     string  `json:"task"`
	Title    string  `json:"title"`
	Solved   string  `json:"solved"`
	HitRatio float64 `json:"hitratio"`
}

func printResp(resp *http.Response) {
	body := &bytes.Buffer{}
	_, err := body.ReadFrom(resp.Body)
	check(err)

	fmt.Println(body)
	fmt.Println(resp.Header)
	fmt.Println(resp.StatusCode)
}

func newfileUploadRequestPost(uri string, body *bytes.Buffer, cookie string, contentType string) (*http.Request, error) {

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Cookie", cookie)

	return req, err
}

func newfileUploadRequest(uri string, params map[string]string, path string, cookie string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	filename := filepath.Base(path)
	if filepath.Ext(path) == ".java" {
		filename = "Solution.java"
	}
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return newfileUploadRequestPost(uri, body, cookie, writer.FormDataContentType())
}

func submitRequest(opts map[string]string, filename string, cookie string) string {
	request, err := newfileUploadRequest("https://cses.fi/course/send.php", opts, filename, cookie)
	check(err)

	client := &http.Client{}
	resp, err := client.Do(request)
	check(err)

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	check(err)

	link, _ := doc.Find("a.current:nth-child(12)").Attr("href")

	return link
}

func loginRequest(params string, cookie string) string {

	body := strings.NewReader(params)
	req, err := http.NewRequest("POST", "https://cses.fi/login", body)
	check(err)

	req.Header.Set("Cookie", cookie)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://cses.fi")
	req.Header.Set("Referer", "https://cses.fi/login")

	resp, err := http.DefaultClient.Do(req)
	check(err)
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	check(err)

	return doc.Find(".account").Contents().Text()
}

func listRequest(cookie string) []*Problem {
	req, err := http.NewRequest("GET", "https://cses.fi/problemset/list", nil)
	check(err)

	req.Header.Set("Cookie", cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	check(err)
	defer resp.Body.Close()

	var arr = []*Problem{}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
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
		arr = append(arr, &Problem{taskNumber, title, solved, percent})
	})
	return arr
}

func printResultRequest(link string, cookie string) (string, string, string) {
	req, err := http.NewRequest("GET", "https://cses.fi"+link, nil)
	check(err)

	req.Header.Set("Cookie", cookie)

	resp, err := http.DefaultClient.Do(req)
	check(err)
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	check(err)

	status := doc.Find("#status").Text()
	text := doc.Find(".summary-table > tbody:nth-child(2)").Contents().Text()

	verdict := doc.Find(".summary-table > tbody:nth-child(2) > tr:nth-child(6) > td:nth-child(2) > span:nth-child(1)").Contents().Text()

	return status, text, verdict
}

func downloadTask(task string) string {

	req, err := http.NewRequest("GET", "https://cses.fi/problemset/task/"+task, nil)
	check(err)
	resp, err := http.DefaultClient.Do(req)
	check(err)
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	check(err)

	out, err := doc.Find(".content").Html()
	check(err)

	title, err := doc.Find("title").Html()
	check(err)

	return title + out
}

func newCookieCsrf() (string, string) {
	req, err := http.NewRequest("GET", "https://cses.fi/login/", nil)
	check(err)

	resp, err := http.DefaultClient.Do(req)
	check(err)

	defer resp.Body.Close()

	cookie := resp.Header.Get("Set-Cookie")

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	check(err)
	csrf, _ := doc.Find(".content > form:nth-child(2) > input:nth-child(1)").Attr("value")

	return cookie, csrf
}
