package main

import "encoding/base64"

func init() {
	funcMap["base64Decode"] = &tmplFuncStruct{
		short: "Decodes a base64 string.",
		examples: []string{
			`{{ "SGVsbG8gV09STEQh" | %s }}`,
		},
		fn: func(str string) (string, error) {
			data, err := base64.StdEncoding.DecodeString(str)
			if err != nil {
				return "", err
			}
			return string(data), nil
		},
	}
}
