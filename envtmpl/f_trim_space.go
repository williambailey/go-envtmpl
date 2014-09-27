package main

import "strings"

func init() {
	funcMap["trimSpace"] = &tmplFuncStruct{
		short: "Remove all leading and trailing white space.",
		examples: []string{
			`{{ " \t\n foo bar \t\n " | %s }}`,
		},
		fn: strings.TrimSpace,
	}
}
