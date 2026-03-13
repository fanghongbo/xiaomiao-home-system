package utils

import "strconv"

func RemoveDuplicate[T comparable](slice []T) []T {
	keys := make(map[T]bool)
	list := []T{}
	for _, item := range slice {
		if _, value := keys[item]; !value {
			keys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func CheckExists[T comparable](slice []T, item T) bool {
	for _, entry := range slice {
		if entry == item {
			return true
		}
	}
	return false
}

func StrToInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}
