package main

import "encoding/base32"

func init() {
	funcMap["base32Encode"] = &tmplFuncStruct{
		short: "Encodes a value to base32.",
		examples: []string{
			`{{ "Hello WORLD!" | %s }}`,
		},
		fn: func(src string) (string, error) {
			data := []byte(src)
			return base32.StdEncoding.EncodeToString(data), nil
		},
	}
}
