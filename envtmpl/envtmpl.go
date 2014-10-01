package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const usageTemplate = `
Usage:
  {{ .cmd }} tmplDir tmplName.tmpl

Parse tmplDir/*.tmpl and renders tmplName.tmpl to
STDOUT using environment variables.

Help:
  {{ .cmd }} -h
`

const helpTemplate = `
# Usage: {{ .cmd }} tmplDir tmplName.tmpl

Parse **tmplDir/*.tmpl** and renders **tmplName.tmpl** to STDOUT using
environment variables.

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
{{ end }}
{{ end }}
`

const exitOk = 0
const exitUsage = 1
const exitTemplateParseError = 2
const exitTemplateExecutionError = 3

var (
	funcMap = newTmplFuncMap()
)

func main() {
	os.Exit(new(os.Environ(), os.Args, os.Stdout, os.Stderr).main())
}

func new(
	env []string,
	args []string,
	stdout io.Writer,
	stderr io.Writer,
) *envtmpl {
	f := flag.NewFlagSet(filepath.Base(args[0]), flag.ExitOnError)
	app := &envtmpl{
		stdout:   stdout,
		stderr:   stderr,
		env:      env,
		cmd:      filepath.Base(args[0]),
		flag:     f,
		flagHelp: f.Bool("h", false, "Display help information, including function list."),
	}
	app.flag.Usage = app.usage
	app.flag.Parse(args[1:])
	return app
}

type envtmpl struct {
	env      []string
	stdout   io.Writer
	stderr   io.Writer
	cmd      string
	flag     *flag.FlagSet
	flagHelp *bool
}

func (app *envtmpl) main() int {
	if *app.flagHelp {
		app.helpUsage()
		return exitUsage
	}
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
	var u bytes.Buffer
	usageTemplate.Execute(&u, map[string]interface{}{
		"cmd": filepath.Base(app.cmd),
	})
	fmt.Fprintf(app.stderr, "%s\n\n", bytes.TrimSpace(u.Bytes()))
}

func (app *envtmpl) helpUsage() {
	usageTemplate := template.Must(
		template.New("").Funcs(funcMap.funcs()).Parse(helpTemplate),
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
	fmt.Fprintf(app.stderr, "%s\n\n", bytes.TrimSpace(u.Bytes()))
}

type tmplFuncMap map[string]tmplFunc

func newTmplFuncMap() tmplFuncMap {
	return make(tmplFuncMap)
}

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
