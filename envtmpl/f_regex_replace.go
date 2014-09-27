package main

import "regexp"

func init() {
	funcMap["regexReplace"] = &tmplFuncStruct{
		short: "Replace values using a regular expression.",
		examples: []string{
			`{{ "this is something" | %s "(this) is " "[$1] was " }}`,
		},
		fn: func(search, replace, src string) (string, error) {
			re, err := regexp.Compile(search)
			if err != nil {
				return "", err
			}
			return re.ReplaceAllString(src, replace), nil
		},
	}
}
