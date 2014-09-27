package main

import (
	"bytes"
	"strings"
)

func init() {
	funcMap["linePrefix"] = &tmplFuncStruct{
		short: "Prefix each line.",
		examples: []string{
			`{{ "line1\nline2\nline3" | %s "- " }}`,
		},
		fn: func(prefix, data string) string {
			var b bytes.Buffer
			for _, v := range strings.SplitAfter(data, "\n") {
				b.WriteString(prefix)
				b.WriteString(v)
			}
			return b.String()
		},
	}
}
