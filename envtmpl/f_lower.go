package main

import "strings"

func init() {
	funcMap["lower"] = &tmplFuncStruct{
		short: "Convert to lower case.",
		examples: []string{
			`{{ "Hello WORLD!" | %s }}`,
		},
		fn: strings.ToLower,
	}
}
