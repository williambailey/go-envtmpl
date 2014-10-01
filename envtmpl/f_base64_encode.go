package main

import "encoding/base64"

func init() {
	funcMap["base64Encode"] = &tmplFuncStruct{
		short: "Encodes a value to base64.",
		examples: []string{
			`{{ "Hello WORLD!" | %s }}`,
		},
		fn: func(src string) (string, error) {
			data := []byte(src)
			return base64.StdEncoding.EncodeToString(data), nil
		},
	}
}
