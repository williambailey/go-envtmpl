package main

import "strings"

func init() {
	funcMap["trimSuffix"] = &tmplFuncStruct{
		short: "Remove trailing suffix. If the string doesn't end with the suffix then it's unchanged.",
		examples: []string{
			`{{ "foo.bar" | %s ".bar" }}`,
			`{{ "foo.bar" | %s ".baz" }}`,
		},
		fn: func(prefix, s string) string {
			return strings.TrimSuffix(s, prefix)
		},
	}
}
