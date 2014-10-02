package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"code.google.com/p/go.text/unicode/norm"
)

func init() {
	allowed := []*unicode.RangeTable{unicode.Letter, unicode.Number}
	dashes := regexp.MustCompile("-+")
	funcMap["slugify"] = &tmplFuncStruct{
		short: "Transform text to a slugified version. You can optionally specify a unicode normalization rule (NFC, NFD, NFKC, NFKD). Default is NFC.",
		examples: []string{
			`{{ "Hello WORLD!" | %s }}`,
			`{{ "Hello W/O-R_L~D!" | %s }}`,
			` NFC: {{ "Hello áçćèńtš!" | %[1]s "NFC" }}
 NFD: {{ "Hello áçćèńtš!" | %[1]s "NFD" }}
NFKC: {{ "Hello áçćèńtš!" | %[1]s "NFKC" }}
NFKD: {{ "Hello áçćèńtš!" | %[1]s "NFKD" }}`,
		},
		fn: func(in ...string) (string, error) {
			var src, mode string
			switch len(in) {
			case 1:
				src = in[0]
			case 2:
				src = in[1]
				mode = in[0]
			default:
				return "", errors.New("Expecting 1 or 2 arguments.")
			}
			switch strings.ToUpper(mode) {
			case "NFC", "C", "":
				src = norm.NFC.String(src)
			case "NFD", "FD":
				src = norm.NFD.String(src)
			case "NFKC", "KC":
				src = norm.NFKC.String(src)
			case "NFKD", "KD":
				src = norm.NFKD.String(src)
			default:
				return "", fmt.Errorf("Unknown normalisation '%s'.", mode)
			}
			runes := make([]rune, len(src))
			for _, r := range src {
				if unicode.IsOneOf(allowed, r) || r == '-' || r == '_' || r == '~' {
					runes = append(runes, r)
				} else if unicode.IsSpace(r) {
					runes = append(runes, '-')
				}
			}
			return strings.ToLower(
				dashes.ReplaceAllString(string(runes), "-"),
			), nil
		},
	}
}
