package main

import (
	"bytes"
	"encoding/json"
)

func init() {
	funcMap["jsonEncode"] = &tmplFuncStruct{
		short: "Encodes a value to JSON.",
		examples: []string{
			`{{ "Hello\n<WORLD>!" | %s }}`,
			`{{ split "Hello\n<WORLD>!" "\n" | %s }}`,
		},
		fn: func(src interface{}) (string, error) {
			b, err := json.Marshal(src)
			if err != nil {
				return "", err
			}
			var o bytes.Buffer
			json.Indent(&o, b, "", "  ")
			return o.String(), nil
		},
	}
}
