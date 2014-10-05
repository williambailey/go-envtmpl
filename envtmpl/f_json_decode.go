package main

import "encoding/json"

func init() {
	funcMap["jsonDecode"] = &tmplFuncStruct{
		short: "Decodes a JSON string.",
		examples: []string{
			`{{ $j := "{\"foo\":\"bar\"}" | %s }}Foo is {{ $j.foo }}`,
		},
		fn: func(src string) (interface{}, error) {
			var o interface{}
			err := json.Unmarshal([]byte(src), &o)
			if err != nil {
				return nil, err
			}
			return o, nil
		},
	}
}
