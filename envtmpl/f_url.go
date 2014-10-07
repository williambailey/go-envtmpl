package main

import (
	"errors"
	"net/url"
)

func init() {
	funcMap["url"] = &tmplFuncStruct{
		short: "Parse a URL from a string. You can optionally provide a base URL that will be used as the context for processing.",
		examples: []string{
			`{{define "u"}}Url: {{ . }}
IsAbs: {{ .IsAbs }}
Scheme: "{{ .Scheme }}"
Opaque: "{{ .Opaque }}"
UserPassword: "{{ .UserPassword }}"
Username: "{{ .Username }}"
HasPassword: {{ .HasPassword }}
Password: "{{ .Password }}"
Host: "{{ .Host }}"
Path: "{{ .Path }}"
RawQuery: "{{ .RawQuery }}"
Query: {{ .Query | jsonEncode }}
Fragment: "{{ .Fragment }}"
RequestURI: "{{ .RequestURI }}"
{{end}}
{{ template "u" "scheme://username:password@domain:port/path?query=string#fragment_id" | %[1]s }}
{{ template "u" "scheme://username:@domain?a=1&b=2&a=11" | %[1]s }}
{{ template "u" "scheme:opaque?query=string#fragment_id" | %[1]s }}
{{ template "u" "../bar/baz" | %[1]s }}
{{ template "u" "../bar/baz" | %[1]s "scheme://domain/foo/qux/" }}`,
		},
		fn: func(in ...string) (*tURL, error) {
			var (
				u   *url.URL
				err error
			)
			switch len(in) {
			case 1:
				u, err = url.Parse(in[0])
			case 2:
				u, err = url.Parse(in[0])
				if err == nil {
					u, err = u.Parse(in[1])
				}
			default:
				return &tURL{}, errors.New("Expecting 1 or 2 arguments.")
			}
			if err != nil {
				return &tURL{}, err
			}
			return &tURL{u}, nil
		},
	}
}

type tURL struct {
	*url.URL
}

func (u *tURL) Username() string {
	if u.User == nil {
		return ""
	}
	return u.User.Username()
}

func (u *tURL) HasPassword() bool {
	if u.User == nil {
		return false
	}
	_, b := u.User.Password()
	return b
}

func (u *tURL) Password() string {
	if u.User == nil {
		return ""
	}
	p, _ := u.User.Password()
	return p
}

func (u *tURL) UserPassword() string {
	if u.User == nil {
		return ""
	}
	return u.User.String()
}
