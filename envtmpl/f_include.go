package main

import (
	"bytes"
	"text/template"
)

func init() {
	funcMap["include"] = &tmplFuncFactoryStruct{
		short:    "Include a template. Differs from the template keyword in that it can accept the template name as part of a pipeline.",
		examples: []string{`{{define "ex"}}FOO is {{ .FOO }}{{end}}{{ $t := "ex" }}>>{{ %s $t . }}<<`},
		fn: func(t *template.Template) interface{} {
			return func(template string, data interface{}) (string, error) {
				var b bytes.Buffer
				err := t.ExecuteTemplate(&b, template, data)
				return b.String(), err
			}
		},
	}
}
