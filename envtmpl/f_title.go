package main

import "strings"

func init() {
	funcMap["title"] = &tmplFuncStruct{
		short: "Convert to title case.",
		examples: []string{
			`{{ "foo BAR bAz" | %s }}`,
		},
		fn: strings.Title,
	}
}
