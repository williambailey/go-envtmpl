package main

import "strings"

func init() {
	funcMap["split"] = &tmplFuncStruct{
		short: "Split in a string substrings using another string.",
		examples: []string{
			`{{ range $k, $v := %s "foo BAR bAz" " " }}{{ $k }}={{ $v }} {{ end }}`,
		},
		fn: strings.Split,
	}
}
