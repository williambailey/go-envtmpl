package main

import "strings"

func init() {
	funcMap["trimPrefix"] = &tmplFuncStruct{
		short: "Remove leading prefix. If the string doesn't start with the prefix then it's unchanged.",
		examples: []string{
			`{{ "foo.bar" | %s "foo." }}`,
			`{{ "foo.bar" | %s "baz." }}`,
		},
		fn: func(prefix, s string) string {
			return strings.TrimPrefix(s, prefix)
		},
	}
}
