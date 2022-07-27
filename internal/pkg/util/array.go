package util

import "strings"

func ArrayIsContains(items []string, item string) bool {
	for _, eachItem := range items {
		if strings.Contains(item, eachItem) {
			return true
		}
	}
	return false
}
func InArray(items []string, item string) bool {
	for _, item1 := range items {
		if item1 == item {
			return true
		}
	}
	return false
}
