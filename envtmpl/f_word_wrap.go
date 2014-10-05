package main

import "strings"

func init() {
	funcMap["wordWrap"] = &tmplFuncStruct{
		short: "Wraps text to a given number of runes. Any existing white space is lost in the transformation.",
		examples: []string{
			`{{ "The quick brown fox jumps over the lazy dog." | %s 19 }}`,
			`{{ "\t  The quick\nbrown fox jumps over the\n\t\tlazy dog." | %s 19 }}`,
			`{{ "Γαζέες καὶ μυρτιὲς δὲν θὰ βρῶ πιὰ στὸ χρυσαφὶ ξέφωτο" | %s 19 }}`,
		},
		fn: func(n int, src string) string {
			var o []rune
			var rl int
			l := n
			for i, w := range strings.Fields(src) {
				r := []rune(w)
				rl = len(r)
				if rl+1 > l {
					o = append(o, '\n')
					l = n - rl
				} else {
					if i > 0 {
						o = append(o, ' ')
					}
					l -= (rl + 1)
				}
				o = append(o, r...)
			}
			return string(o)
		},
	}
}
