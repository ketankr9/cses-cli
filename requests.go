package main

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// func printResp(resp *http.Response) {
// 	body := &bytes.Buffer{}
// 	_, err := body.ReadFrom(resp.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(body)
// 	fmt.Println(resp.Header)
// 	fmt.Println(resp.StatusCode)
// }

func newfileUploadRequestPost(uri string, body *bytes.Buffer, cookie string, contentType string) (*http.Request, error) {

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Cookie", cookie)

	return req, err
}

func submitRequest(opts map[string]string, filename string, cookie string) string {
	request, err := newfileUploadRequest("https://cses.fi/course/send.php", opts, "file", filename, cookie)
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

func listRequest(cookie string) io.ReadCloser {
	req, err := http.NewRequest("GET", "https://cses.fi/problemset/list", nil)
	check(err)

	req.Header.Set("Cookie", cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	check(err)

	// defer resp.Body.Close()

	return resp.Body
}

func printResultRequest(link string, cookie string) (string, string) {
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

	return status, text
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
