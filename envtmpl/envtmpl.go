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

const Name = "envtmpl"
const Version = "0.1.0"

const usageTemplate = `
Usage:
  {{ .cmd }} tmplDir tmplName.tmpl
  {{ .cmd }} tmplDir/tmplName.tmpl
  {{ .cmd }} -

Parse tmplDir/*.tmpl and renders tmplName.tmpl to
STDOUT using environment variables. If a dash is
provided then the template is read from STDIN.

Version:
  {{ .version }}

Help:
  {{ .cmd }} -h
`

const helpTemplate = `
# Usage:

#### {{ .usage1 }}

Parse **tmplDir/*.tmpl** and renders **tmplName.tmpl** to STDOUT using
environment variables.

#### {{ .usage2 }}

Parse **tmplDir/*.tmpl** and renders **tmplName.tmpl** to STDOUT using
environment variables.

#### {{ .usageStdin }}

Read template from STDIN and render to STDOUT using environment variables.

### Exit codes

* 0 - OK.
* 1 - Usage.
* 2 - Template parse error.
* 3 - Template execution error.

### Template Syntax.

See http://golang.org/pkg/text/template/ for template syntax.

See http://code.google.com/p/re2/wiki/Syntax for regular expression syntax.

## Template Functions

In addition to the [actions](http://golang.org/pkg/text/template/#hdr-Actions)
and [functions](http://golang.org/pkg/text/template/#hdr-Functions) provided by
the core [template engine](http://golang.org/pkg/text/template/#pkg-overview),
{{ .cmd }} provides the following functions for use in your templates:

{{ range .funcs }}* [{{ .Name }}](#{{ .Name | slugify }}) - {{ index (split .Short ".") 0 }}.
{{ end }}
{{ range .funcs }}{{ template "funcHelp" . }}{{ end }}
`

const funcHelpTemplate = `### {{ .Name }}

{{ .Short | wordWrap 80 }}
{{ range $ex, $out := .Example }}
Template:

{{ $ex | linePrefix "    " }}

Output:

{{ $out | linePrefix "    " }}
{{ end }}
`

const exitOk = 0
const exitUsage = 1
const exitTemplateParseError = 2
const exitTemplateExecutionError = 3

var (
	funcMap         = newTmplFuncMap()
	funcHelpExample = false
)

func main() {
	os.Exit(new(os.Environ(), os.Args, os.Stdin, os.Stdout, os.Stderr).main())
}

func new(
	env []string,
	args []string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) *envtmpl {
	f := flag.NewFlagSet(filepath.Base(args[0]), flag.ExitOnError)
	app := &envtmpl{
		stdin:    stdin,
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
	stdin    io.Reader
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
	args := app.flag.Args()
	var tmplDir string
	var tmplName string
	tmplData := make(map[string]string)
	switch len(args) {
	case 1:
		if args[0] == "-" {
			tmplDir = "-"
			tmplName = "stdin"
		} else {
			tmplDir = filepath.Dir(args[0])
			tmplName = filepath.Base(args[0])
		}
	case 2:
		tmplDir = args[0]
		tmplName = args[1]
	default:
		app.flag.Usage()
		return exitUsage
	}
	tmpl := template.New(
		fmt.Sprintf("%s [%s]", app.cmd, tmplDir),
	).Funcs(funcMap.funcs())

	var err error
	if tmplDir == "-" && tmplName == "stdin" {
		var b bytes.Buffer
		b.ReadFrom(app.stdin)
		_, err = tmpl.New("stdin").Parse(b.String())
	} else {
		_, err = tmpl.ParseGlob(filepath.Join(tmplDir, "*.tmpl"))
	}
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
	t := template.Must(
		template.New("usage").Funcs(funcMap.funcs()).Parse(usageTemplate),
	)
	var u bytes.Buffer
	t.Execute(&u, map[string]interface{}{
		"version": Version,
		"cmd":     filepath.Base(app.cmd),
	})
	fmt.Fprintf(app.stderr, "%s\n", bytes.TrimSpace(u.Bytes()))
}

func (app *envtmpl) helpUsage() {
	funcHelpExample = true
	defer func() { funcHelpExample = false }()
	t := template.Must(
		template.New("help").Funcs(funcMap.funcs()).Parse(helpTemplate),
	)
	template.Must(t.New("funcHelp").Parse(funcHelpTemplate))

	type funcData struct {
		Name    string
		Short   string
		Example map[string]string
	}
	var u bytes.Buffer
	funcs := make(map[string]funcData)
	for n, fn := range funcMap {
		fd := funcData{
			Name:    n,
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
	cmd := filepath.Base(app.cmd)
	err := t.Execute(&u, map[string]interface{}{
		"cmd":        cmd,
		"funcs":      funcs,
		"usage1":     cmd + " tmplDir tmplName.tmpl",
		"usage2":     cmd + " tmplDir/tmplName.tmpl",
		"usageStdin": cmd + " -",
	})
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(app.stderr, "%s\n", bytes.TrimSpace(u.Bytes()))
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
