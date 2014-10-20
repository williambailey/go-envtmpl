package main

import "net/url"

func init() {
	funcMap["urlUnescape"] = &tmplFuncStruct{
		short: "Unescapes a URL query string value.",
		examples: []string{
			`{{ "Hello+World%%21" | %s }}`,
		},
		fn: url.QueryUnescape,
	}
}
