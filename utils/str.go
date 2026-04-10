package utils

import "unicode/utf8"

func GetUtf8RuneCount(s string) int {
	return utf8.RuneCountInString(s)
}
