package main

import "encoding/hex"

func init() {
	funcMap["hexEncode"] = &tmplFuncStruct{
		short: "Encodes a value to hex.",
		examples: []string{
			`{{ "Hello WORLD!" | %s }}`,
		},
		fn: func(src string) (string, error) {
			data := []byte(src)
			return hex.EncodeToString(data), nil
		},
	}
}
