package main

import "strings"

func init() {
	funcMap["upper"] = &tmplFuncStruct{
		short: "Convert to upper case.",
		examples: []string{
			`{{ "Hello World!" | %s }}`,
		},
		fn: strings.ToUpper,
	}
}
