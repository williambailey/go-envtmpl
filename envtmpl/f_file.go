package main

import "io/ioutil"

func init() {
	funcMap["file"] = &tmplFuncStruct{
		short:    "Read the contents of a file.",
		examples: []string{`{{ %s "../example/hello.txt" }}`},
		fn: func(file string) (string, error) {
			if funcHelpExample {
				return "Hello, 世界", nil
			}
			b, err := ioutil.ReadFile(file)
			return string(b), err
		},
	}
}
