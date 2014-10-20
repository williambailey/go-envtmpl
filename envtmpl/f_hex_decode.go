package main

import "encoding/hex"

func init() {
	funcMap["hexDecode"] = &tmplFuncStruct{
		short: "Decodes a hex string.",
		examples: []string{
			`{{ "48656c6c6f20574f524c4421" | %s }}`,
		},
		fn: func(str string) (string, error) {
			data, err := hex.DecodeString(str)
			if err != nil {
				return "", err
			}
			return string(data), nil
		},
	}
}
