package util

import "strings"

func SetCaseForFirstLetter(s string, upper bool) string {
	if len(s) == 0 {
		return s
	}
	// TODO improve for unicode
	setter := strings.ToLower
	if upper {
		setter = strings.ToUpper
	}
	return setter(s[:1]) + s[1:]
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
