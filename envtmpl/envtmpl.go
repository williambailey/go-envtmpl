package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"code.google.com/p/go-uuid/uuid"
)

const usageTemplate = `
# Usage: {{ .cmd }} tmplDir tmplName.tmpl

Parse **tmplDir/*.tmpl** and renders **tmplName.tmpl** to STDOUT using
the data from environmental variables.

See http://golang.org/pkg/text/template/ for template syntax.

See http://code.google.com/p/re2/wiki/Syntax for regular expression syntax.

### Exit codes

* 0 - OK.
* 1 - Usage.
* 2 - Template parse error.
* 3 - Template execution error.

## Template Functions

{{ range $k, $v := .funcs }}### {{ $k }}

{{ $v.Short }}
{{ range $ex, $out := $v.Example }}
Template:

{{ $ex | linePrefix "    " }}

Output:

{{ $out | linePrefix "    " }}

{{ end }}{{ end }}
`

const exitOk = 0
const exitUsage = 1
const exitTemplateParseError = 2
const exitTemplateExecutionError = 3

var (
	funcMap tmplFuncMap
)

func init() {
	funcMap = make(tmplFuncMap)
	initFuncMap()
}

func main() {
	os.Exit(new(os.Environ(), os.Args, os.Stdout, os.Stderr).main())
}

func new(
	env []string,
	args []string,
	stdout io.Writer,
	stderr io.Writer,
) *envtmpl {
	app := &envtmpl{
		stdout: stdout,
		stderr: stderr,
		env:    env,
		cmd:    filepath.Base(args[0]),
		flag:   flag.NewFlagSet(filepath.Base(args[0]), flag.ExitOnError),
	}
	app.flag.Usage = app.usage
	app.flag.Parse(args[1:])
	return app
}

type envtmpl struct {
	env    []string
	stdout io.Writer
	stderr io.Writer
	cmd    string
	flag   *flag.FlagSet
}

func (app *envtmpl) main() int {
	if app.flag.NArg() != 2 {
		app.flag.Usage()
		return exitUsage
	}
	args := app.flag.Args()
	tmplDir := args[0]
	tmplName := args[1]
	tmplData := make(map[string]string)

	tmpl := template.New(
		fmt.Sprintf("%s [%s]", app.cmd, tmplDir),
	).Funcs(funcMap.funcs())

	_, err := tmpl.ParseGlob(filepath.Join(tmplDir, "*.tmpl"))
	if err != nil {
		fmt.Fprintf(app.stderr, "Template parse error: %s\n", err)
		return exitTemplateParseError
	}

	for _, s := range app.env {
		o := strings.Index(s, "=")
		if o <= 0 {
			continue
		}
		tmplData[s[:o]] = s[o+1:]
	}

	err = tmpl.ExecuteTemplate(app.stdout, tmplName, tmplData)
	if err != nil {
		fmt.Fprintf(app.stderr, "Template execution: %s\n", err)
		return exitTemplateExecutionError
	}
	return exitOk
}

func (app *envtmpl) usage() {
	usageTemplate := template.Must(
		template.New("").Funcs(funcMap.funcs()).Parse(usageTemplate),
	)
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
			err := template.Must(
				template.New(n).Funcs(funcMap.funcs()).Parse(e),
			).Execute(&b, struct{}{})
			if err != nil {
				panic(err)
			}
			fd.Example[e] = b.String()
		}
		funcs[n] = fd
	}
	usageTemplate.Execute(&u, map[string]interface{}{
		"cmd":   filepath.Base(app.cmd),
		"funcs": funcs,
	})
	fmt.Fprintf(app.stderr, "%s", bytes.TrimSpace(u.Bytes()))
	fmt.Fprint(app.stderr, "\n\n")
	app.flag.PrintDefaults()
}

type tmplFuncMap map[string]tmplFunc

func (m *tmplFuncMap) funcs() map[string]interface{} {
	funcs := make(template.FuncMap)
	for k, v := range funcMap {
		funcs[k] = v.f()
	}
	return funcs
}

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
	e := make([]string, len(s.examples))
	for k, v := range s.examples {
		e[k] = fmt.Sprintf(v, name)
	}
	return e
}

func (s *tmplFuncStruct) f() interface{} {
	return s.fn
}

func initFuncMap() {
	const fooExample = `{{ "foo BAR bAz" | %s }}`
	funcMap["linePrefix"] = &tmplFuncStruct{
		short: "Prefix each line.",
		examples: []string{
			`{{ "line1\nline2\nline3" | %s "- " }}`,
		},
		fn: func(prefix, data string) string {
			var b bytes.Buffer
			for _, v := range strings.SplitAfter(data, "\n") {
				b.WriteString(prefix)
				b.WriteString(v)
			}
			return b.String()
		},
	}
	funcMap["lower"] = &tmplFuncStruct{
		short: "Convert to lower case.",
		examples: []string{
			fooExample,
		},
		fn: strings.ToLower,
	}
	funcMap["regexReplace"] = &tmplFuncStruct{
		short: "Replace values using a regular expression.",
		examples: []string{
			`{{ "this is something" | %s "(this) is " "[$1] was " }}`,
		},
		fn: func(search, replace, src string) (string, error) {
			re, err := regexp.Compile(search)
			if err != nil {
				return "", err
			}
			return re.ReplaceAllString(src, replace), nil
		},
	}
	funcMap["split"] = &tmplFuncStruct{
		short: "Split in a string substrings using another string.",
		examples: []string{
			`{{ range $k, $v := %s "foo BAR bAz" " " }}{{ $k }}={{ $v }} {{ end }}`,
		},
		fn: strings.Split,
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
