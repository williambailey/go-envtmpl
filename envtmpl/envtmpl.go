package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"code.google.com/p/go-uuid/uuid"
)

const exitOk = 0
const exitUsage = 1
const exitTemplateParseError = 2
const exitTemplateExecutionError = 3

var (
	funcMap   tmplFuncMap
	tmplFuncs template.FuncMap
)

func init() {
	funcMap = make(tmplFuncMap)
	initFuncMap()
	tmplFuncs = make(template.FuncMap)
	for k, v := range funcMap {
		tmplFuncs[k] = v.f()
	}
}

func main() {
	os.Exit(realMain())
}

func realMain() int {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		return exitUsage
	}
	args := flag.Args()
	tmplDir := args[0]
	tmplName := args[1]
	tmplData := make(map[string]string)

	tmpl := template.New(
		fmt.Sprintf("%s [%s]", filepath.Base(os.Args[0]), tmplDir),
	).Funcs(tmplFuncs)

	_, err := tmpl.ParseGlob(filepath.Join(tmplDir, "*.tmpl"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Template parse error: %s\n", err)
		return exitTemplateParseError
	}

	for _, s := range os.Environ() {
		o := strings.Index(s, "=")
		if o <= 0 {
			continue
		}
		tmplData[s[:o]] = s[o+1:]
	}

	err = tmpl.ExecuteTemplate(os.Stdout, tmplName, tmplData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Template execution: %s\n", err)
		return exitTemplateExecutionError
	}
	return exitOk
}

func usage() {
	usageTemplate := template.Must(template.New("").Parse(`
Usage: {{ .cmd }} tmplDir tmplName

Parse tmplDir/*.tmpl and renders tmplName to stdout using
the data from environmental variables.

See http://golang.org/pkg/text/template/ for template syntax.

Exit codes:

0 OK.
1 Usage.
2 Template parse error.
3 Template execution error.

Template Functions:

{{ range $k, $v := .funcs }}{{ $k }} :: {{ $v.Short }}
{{ range $ex, $out := $v.Example }} * {{ $ex }}
{{ $out }}
{{ end }}
{{ end }}

`))
	type funcData struct {
		Short   string
		Example map[string]string
	}
	var u bytes.Buffer
	funcs := make(map[string]funcData)
	for n, fn := range funcMap {
		fd := funcData{
			Short:   fn.shortUsage(),
			Example: make(map[string]string),
		}
		for _, e := range fn.example(n) {
			var b bytes.Buffer
			ex := fmt.Sprintf(e, n)
			template.Must(
				template.New(n).Funcs(tmplFuncs).Parse(ex),
			).Execute(&b, struct{}{})
			fd.Example[ex] = b.String()
		}
		funcs[n] = fd
	}
	usageTemplate.Execute(&u, map[string]interface{}{
		"cmd":   filepath.Base(os.Args[0]),
		"funcs": funcs,
	})
	fmt.Fprintf(os.Stderr, "%s", bytes.TrimSpace(u.Bytes()))
	fmt.Fprint(os.Stderr, "\n\n")
	flag.PrintDefaults()
}

type tmplFuncMap map[string]tmplFunc

type tmplFunc interface {
	shortUsage() string
	example(name string) []string
	f() interface{}
}

type tmplFuncStruct struct {
	short    string
	examples []string
	fn       interface{}
}

func (s *tmplFuncStruct) shortUsage() string {
	return s.short
}

func (s *tmplFuncStruct) example(name string) []string {
	return s.examples
}

func (s *tmplFuncStruct) f() interface{} {
	return s.fn
}

func initFuncMap() {
	const fooExample = `{{ "foo BAR bAz" | %s }}`
	funcMap["lower"] = &tmplFuncStruct{
		short: "Convert to lower case.",
		examples: []string{
			fooExample,
		},
		fn: strings.ToLower,
	}
	funcMap["title"] = &tmplFuncStruct{
		short: "Convert to title case.",
		examples: []string{
			fooExample,
		},
		fn: strings.Title,
	}
	funcMap["upper"] = &tmplFuncStruct{
		short: "Convert to upper case.",
		examples: []string{
			fooExample,
		},
		fn: strings.ToUpper,
	}
	funcMap["uuid"] = &tmplFuncStruct{
		short: "Create a random (v4) UUID.",
		examples: []string{
			`{{ %s }}`,
		},
		fn: uuid.New,
	}
}
