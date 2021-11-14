package pkg

import (
	"strings"
	"unicode"
)

func Slugify(s string) string {
	return replaceSpaces(strings.ToLower(s))
}

func replaceSpaces(s string) string {
	runes := make([]rune, len(s))
	for i, v := range s {
		if unicode.IsSpace(v) {
			v = '-'
		}
		runes[i] = v
	}
	return string(runes)
}
