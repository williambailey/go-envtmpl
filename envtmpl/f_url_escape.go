package main

import "net/url"

func init() {
	funcMap["urlEscape"] = &tmplFuncStruct{
		short: "Escapes the string so it can be safely placed inside a URL query.",
		examples: []string{
			`{{ "Hello World!" | %s }}`,
		},
		fn: url.QueryEscape,
	}
}
