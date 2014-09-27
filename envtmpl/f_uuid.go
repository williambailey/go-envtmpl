package main

import "code.google.com/p/go-uuid/uuid"

func init() {
	funcMap["uuid"] = &tmplFuncStruct{
		short: "Create a random (v4) UUID.",
		examples: []string{
			`{{ %s }}`,
		},
		fn: uuid.New,
	}
}
