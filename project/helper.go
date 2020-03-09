package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func validGithubConfig(opts *githubConfig) bool {
	return opts.Token != "" && opts.SourceRepo != "" && opts.AuthorName != "" && opts.AuthorEmail != ""
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

func updateConfig(sess *Session) {
	out, err := json.MarshalIndent(sess, "", "    ")
	check(err)
	cacheSet("config.json", string(out), sess.Root)
}

func updateIfNew(scanner *bufio.Scanner, src *string, text string) {
	if *src == "" {
		fmt.Print(text + ": ")
	} else {
		fmt.Print(text + "(" + *src + "): ")
	}
	scanner.Scan()
	val := scanner.Text()
	if val != "" {
		*src = val
	}
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

func getTemplate(ext string) string {
	content, err := ioutil.ReadFile("template"+ext)
	if err != nil {
		return ""
	}
	return string(content)
}

func getTaskFromCache(task string, root string) string {
	path := filepath.Join(root, task+".task.html")

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

var extLangMap = map[string]string{
	".java": "Java",
	".js": "Node.js",
	".py": "Python3",
	".cpp": "C++",
}
var extOptionMap = map[string]string{
	".java": "",
	".js": "",
	".py": "CPython3",
	".cpp": "C++17",
}
var langExtMap = map[string]string {
	"java": ".java",
	"javascript": ".js",
	"python": ".py",
	"cpp": ".cpp",
}

var unicodeMap = map[string]string{
	"\\alpha":              "\u03B1",
	"\\beta":               "\u03B2",
	"\\gamma":              "\u03B3",
	"\\delta":              "\u03B4",
	"\\epsilon":            "\u03F5",
	"\\zeta":               "\u03B6",
	"\\eta":                "\u03B7",
	"\\theta":              "\u03B8",
	"\\iota":               "\u03B9",
	"\\kappa":              "\u03BA",
	"\\lambda":             "\u03BB",
	"\\mu":                 "\u03BC",
	"\\nu":                 "\u03BD",
	"\\xi":                 "\u03BE",
	"\\omicron":            "\u03BF",
	"\\pi":                 "\u03C0",
	"\\rho":                "\u03C1",
	"\\sigma":              "\u03C3",
	"\\tau":                "\u03C4",
	"\\upsilon":            "\u03C5",
	"\\phi":                "\u03D5",
	"\\chi":                "\u03C7",
	"\\psi":                "\u03C8",
	"\\omega":              "\u03C9",
	"\\varepsilon":         "\u03B5",
	"\\vartheta":           "\u03D1",
	"\\varpi":              "\u03D6",
	"\\varrho":             "\u03F1",
	"\\varsigma":           "\u03C2",
	"\\varphi":             "\u03C6",
	"\\S":                  "\u00A7",
	"\\aleph":              "\u2135",
	"\\hbar":               "\u210F",
	"\\imath":              "\u0131",
	"\\jmath":              "\u0237",
	"\\ell":                "\u2113",
	"\\wp":                 "\u2118",
	"\\Re":                 "\u211C",
	"\\Im":                 "\u2111",
	"\\partial":            "\u2202",
	"\\infty":              "\u221E",
	"\\prime":              "\u2032",
	"\\emptyset":           "\u2205",
	"\\nabla":              "\u2207",
	"\\top":                "\u22A4",
	"\\bot":                "\u22A5",
	"\\angle":              "\u2220",
	"\\triangle":           "\u25B3",
	"\\backslash":          "\u2216",
	"\\forall":             "\u2200",
	"\\exists":             "\u2203",
	"\\neg":                "\u00AC",
	"\\lnot":               "\u00AC",
	"\\flat":               "\u266D",
	"\\natural":            "\u266E",
	"\\sharp":              "\u266F",
	"\\clubsuit":           "\u2663",
	"\\diamondsuit":        "\u2662",
	"\\heartsuit":          "\u2661",
	"\\spadesuit":          "\u2660",
	"\\surd":               "\u221A",
	"\\coprod":             "\u2210",
	"\\bigvee":             "\u22C1",
	"\\bigwedge":           "\u22C0",
	"\\biguplus":           "\u2A04",
	"\\bigcap":             "\u22C2",
	"\\bigcup":             "\u22C3",
	"\\int":                "\u222B",
	"\\intop":              "\u222B",
	"\\iint":               "\u222C",
	"\\iiint":              "\u222D",
	"\\prod":               "\u220F",
	"\\sum":                "\u2211",
	"\\bigotimes":          "\u2A02",
	"\\bigoplus":           "\u2A01",
	"\\bigodot":            "\u2A00",
	"\\oint":               "\u222E",
	"\\bigsqcup":           "\u2A06",
	"\\smallint":           "\u222B",
	"\\triangleleft":       "\u25C3",
	"\\triangleright":      "\u25B9",
	"\\bigtriangleup":      "\u25B3",
	"\\bigtriangledown":    "\u25BD",
	"\\wedge":              "\u2227",
	"\\land":               "\u2227",
	"\\vee":                "\u2228",
	"\\lor":                "\u2228",
	"\\cap":                "\u2229",
	"\\cup":                "\u222A",
	"\\ddagger":            "\u2021",
	"\\dagger":             "\u2020",
	"\\sqcap":              "\u2293",
	"\\sqcup":              "\u2294",
	"\\uplus":              "\u228E",
	"\\amalg":              "\u2A3F",
	"\\diamond":            "\u22C4",
	"\\bullet":             "\u2219",
	"\\wr":                 "\u2240",
	"\\div":                "\u00F7",
	"\\mp":                 "\u2213",
	"\\pm":                 "\u00B1",
	"\\circ":               "\u2218",
	"\\bigcirc":            "\u25EF",
	"\\setminus":           "\u2216",
	"\\cdot":               "\u22C5",
	"\\ast":                "\u2217",
	"\\times":              "\u00D7",
	"\\star":               "\u22C6",
	"\\propto":             "\u221D",
	"\\sqsubseteq":         "\u2291",
	"\\sqsupseteq":         "\u2292",
	"\\parallel":           "\u2225",
	"\\mid":                "\u2223",
	"\\dashv":              "\u22A3",
	"\\vdash":              "\u22A2",
	"\\leq":                "\u2264",
	"\\le":                 "\u2264",
	"\\geq":                "\u2265",
	"\\ge":                 "\u2265",
	"\\lt":                 "\u003C",
	"\\gt":                 "\u003E",
	"\\succ":               "\u227B",
	"\\prec":               "\u227A",
	"\\approx":             "\u2248",
	"\\succeq":             "\u2AB0",
	"\\preceq":             "\u2AAF",
	"\\supset":             "\u2283",
	"\\subset":             "\u2282",
	"\\supseteq":           "\u2287",
	"\\subseteq":           "\u2286",
	"\\in":                 "\u2208",
	"\\ni":                 "\u220B",
	"\\notin":              "\u2209",
	"\\owns":               "\u220B",
	"\\gg":                 "\u226B",
	"\\ll":                 "\u226A",
	"\\sim":                "\u223C",
	"\\simeq":              "\u2243",
	"\\perp":               "\u22A5",
	"\\equiv":              "\u2261",
	"\\asymp":              "\u224D",
	"\\smile":              "\u2323",
	"\\frown":              "\u2322",
	"\\ne":                 "\u2260",
	"\\neq":                "\u2260",
	"\\cong":               "\u2245",
	"\\doteq":              "\u2250",
	"\\bowtie":             "\u22C8",
	"\\models":             "\u22A8",
	"\\notChar":            "\u29F8",
	"\\Leftrightarrow":     "\u21D4",
	"\\Leftarrow":          "\u21D0",
	"\\Rightarrow":         "\u21D2",
	"\\leftrightarrow":     "\u2194",
	"\\leftarrow":          "\u2190",
	"\\gets":               "\u2190",
	"\\rightarrow":         "\u2192",
	"\\to":                 "\u2192",
	"\\mapsto":             "\u21A6",
	"\\leftharpoonup":      "\u21BC",
	"\\leftharpoondown":    "\u21BD",
	"\\rightharpoonup":     "\u21C0",
	"\\rightharpoondown":   "\u21C1",
	"\\nearrow":            "\u2197",
	"\\searrow":            "\u2198",
	"\\nwarrow":            "\u2196",
	"\\swarrow":            "\u2199",
	"\\rightleftharpoons":  "\u21CC",
	"\\hookrightarrow":     "\u21AA",
	"\\hookleftarrow":      "\u21A9",
	"\\longleftarrow":      "\u27F5",
	"\\Longleftarrow":      "\u27F8",
	"\\longrightarrow":     "\u27F6",
	"\\Longrightarrow":     "\u27F9",
	"\\Longleftrightarrow": "\u27FA",
	"\\longleftrightarrow": "\u27F7",
	"\\longmapsto":         "\u27FC",
	"\\ldots":              "\u2026",
	"\\cdots":              "\u22EF",
	"\\vdots":              "\u22EE",
	"\\ddots":              "\u22F1",
	"\\dotsc":              "\u2026",
	"\\dotsb":              "\u22EF",
	"\\dotsm":              "\u22EF",
	"\\dotsi":              "\u22EF",
	"\\dotso":              "\u2026",

	"\\uparrow":     "\u2191",
	"\\downarrow":   "\u2193",
	"\\updownarrow": "\u2195",
	"\\Uparrow":     "\u21D1",
	"\\Downarrow":   "\u21D3",
	"\\Updownarrow": "\u21D5",
	"\\rangle":      "\u27E9",
	"\\langle":      "\u27E8",
	"\\rbrace":      "}",
	"\\lbrace":      "{",
	"\\}":           "}",
	"\\{":           "{",
	"\\rceil":       "\u2309",
	"\\lceil":       "\u2308",
	"\\rfloor":      "\u230B",
	"\\lfloor":      "\u230A",
	"\\lbrack":      "[",
	"\\rbrack":      "]",

	"\\[":    "[",
	"\\]":    "]",
	"\\dots": "...",
}
