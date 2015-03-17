package main

import "text/template"

func init() {
	data := make(map[*template.Template]map[interface{}]interface{})
	getDataStore := func(t *template.Template) map[interface{}]interface{} {
		d, ok := data[t]
		if !ok {
			d = make(map[interface{}]interface{})
			data[t] = d
		}
		return d
	}
	funcMap["get"] = &tmplFuncFactoryStruct{
		short: "Get a value from the global data store.",
		examples: []string{
			`{{ define "t1" }}[{{ %s "key" }}]{{ end }}
{{ set "key" "value" }}{{ template "t1" }}`,
		},
		fn: func(t *template.Template) interface{} {
			d := getDataStore(t)
			return func(k interface{}) interface{} {
				v, ok := d[k]
				if !ok {
					return nil
				}
				return v
			}
		},
	}
	funcMap["set"] = &tmplFuncFactoryStruct{
		short: "Set a value in the global data store.",
		examples: []string{
			`{{ define "t1" }}[{{ get "key" }}]{{ end }}
{{ %s "key" "value" }}{{ template "t1" }}`,
		},
		fn: func(t *template.Template) interface{} {
			d := getDataStore(t)
			return func(k interface{}, v interface{}) string {
				d[k] = v
				return ""
			}
		},
	}
}
