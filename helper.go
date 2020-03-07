package main

import(
	"net/http"
	"os"
	"bytes"
	"bufio"
	"mime/multipart"
	"path/filepath"
	"io"
	"io/ioutil"
	"strings"
	"os/exec"
	"runtime"

)

var unicodeMap = map[string]string{
		"\\le":         "\u2264",
		"\\dots":       "...",
		"\\cdots":      "\u22EF",
		"\\rightarrow": "\u21D2",
		"\\times":      "\u00D7",
		"\\alpha":      "\u03B1",
		"\\beta":       "\u03B2",
		"\\gamma":      "\u03B3",
		"\\delta":      "\u03B4",
		"\\epsilon":    "\u03F5",
		"\\zeta":       "\u03B6",
		"\\eta":        "\u03B7",
		"\\theta":      "\u03B8",
		"\\iota":       "\u03B9",
		"\\kappa":      "\u03BA",
		"\\lambda":     "\u03BB",
		// parallel: '2225',
		// mid: '2223',
		// dashv: '22A3',
		// vdash: '22A2',
		"\\leq": "\u2264",
		// geq: '2265',
		"\\ge": "\u2265",
		"\\lt": "\u003C",
		"\\gt": "\u003E",
		// succ: '227B',
		// prec: '227A',
		// approx: '2248',
		// succeq: '2AB0',
		// preceq: '2AAF',
		"\\supset": "\u2283",
		"\\subset": "\u2282",
		// supseteq: '2287',
		// subseteq: '2286',
		"\\uparrow": "\u2191",
		// '\\downarrow': '2193',
		// '\\updownarrow': '2195',
		// '\\Uparrow': '21D1',
		// '\\Downarrow': '21D3',
		// '\\Updownarrow': '21D5',
		"\\backslash": "\\",
		// '\\rangle': '27E9',
		// '\\langle': '27E8',
		"\\rbrace": "}",
		"\\lbrace": "{",
		"\\}": "}",
		"\\{": "{",
		"\\[": "[",
		"\\]": "]",
		"\\rceil": "\u2309",
		"\\lceil": "\u2308",
		"\\rfloor": "\u230B",
		"\\lfloor": "\u230A",
		"\\lbrack": "[",
		"\\rbrack": "]",
	}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func UserHomeDir() string {
    if runtime.GOOS == "windows" {
        home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
        if home == "" {
            home = os.Getenv("USERPROFILE")
        }
        return home
    }
    return os.Getenv("HOME")
}

func cacheSet(filename string, data string, root string) {
	f, err := os.Create(filepath.Join(root, filename))
	check(err)

	w := bufio.NewWriter(f)
	w.WriteString(data)

	w.Flush()
	defer f.Close()
}

func cacheGet(filename string, root string) ([]byte, bool) {
	path := filepath.Join(root, filename)

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, false
	}
	content, err := ioutil.ReadFile(path)
	check(err)
	return content, true
}

func getTaskFromCache(task string, root string) string {
	path := filepath.Join(root, task + ".task.html")

	output, err := exec.Command("bash", "-c", "lynx -dump "+path).Output()
	check(err)
	data := string(output)

	for k, v := range unicodeMap {
		data = strings.Replace(data, k, v, -1)
	}

	return data
}

func writeCodeFile(filename string, text string, template string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}

	f, err := os.Create(filename)
	if err != nil {
		return false
	}
	_, err = f.WriteString("/*\n" + text + "*/\n" + template)
	if err != nil {
		f.Close()
		return false
	}

	err = f.Close()
	if err != nil {
		return false
	}
	return true
}

func newfileUploadRequest(uri string, params map[string]string, paramName, path string, cookie string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
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

func fileMeta(filename string) (string, string, string, bool) {
	parts := strings.Split(filepath.Base(filename), ".")
	lang := parts[2]
	option := ""
	if lang == "cpp" {
		lang = "C++"
		option = "C++17"
	}
	return parts[0], lang, option, true
}