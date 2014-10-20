package main

import "encoding/base32"

func init() {
	funcMap["base32Decode"] = &tmplFuncStruct{
		short: "Decodes a base32 string.",
		examples: []string{
			`{{ "JBSWY3DPEBLU6USMIQQQ====" | %s }}`,
		},
		fn: func(str string) (string, error) {
			data, err := base32.StdEncoding.DecodeString(str)
			if err != nil {
				return "", err
			}
			return string(data), nil
		},
	}
}
