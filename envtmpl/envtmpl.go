package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

const exitOk = 0
const exitUsage = 1
const exitTemplateParseError = 2
const exitTemplateExecutionError = 3

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

	tmpl, err := template.ParseGlob(filepath.Join(tmplDir, "*.tmpl"))
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

	err = tmpl.ExecuteTemplate(os.Stdout, tmplName, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Template execution: %s\n", err)
		return exitTemplateExecutionError
	}
	return exitOk
}

func usage() {
	cmd := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, strings.TrimSpace(`
Usage: %s tmplDir tmplName

  Parse tmplDir/*.tmpl and renders tmplName to stdout using
  the data from environmental variables.

  See http://golang.org/pkg/text/template/ for template syntax.

Exit codes:

  0 OK.
  1 Usage.
  2 Template parse error.
  3 Template execution error.

`)+"\n\n", cmd)
	flag.PrintDefaults()
}
