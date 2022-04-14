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
