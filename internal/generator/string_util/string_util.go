package string_util

import (
	"regexp"
	"strings"
)

// Camel case -> snake case
func CamelToSnake(s string) string {
	s = regexp.MustCompile("([^A-Z])([A-Z])").ReplaceAllString(s, "${1}_${2}")
	s = regexp.MustCompile("([A-Z])([A-Z][a-z])").ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(s)
}

func CountMaxConsecutiveUnderscores(s string, initMax int) int {
	max := initMax
	cur := 0
	anyNotUnderscoreRune := "#"
	for _, r := range s + anyNotUnderscoreRune {
		if r == '_' {
			cur++
		} else {
			if max < cur {
				max = cur
			}
			cur = 0
		}
	}
	return max
}
