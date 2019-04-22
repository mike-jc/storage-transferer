package system

import (
	"fmt"
	"regexp"
)

func RemoveSpecialChars(s string, chars string) string {
	if len(chars) == 0 {
		return s
	}

	regStr := fmt.Sprintf("[%s]+", regexp.QuoteMeta(chars))
	return regexp.MustCompile(regStr).ReplaceAllString(s, "")
}

func TruncateByMaxLen(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if len(s) <= maxLen {
		return s
	}
	return s[0:maxLen]
}
