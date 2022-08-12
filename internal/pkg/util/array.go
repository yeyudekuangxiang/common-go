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

func Intersect(nums1 []string, nums2 []string) []string {
	if len(nums1) > len(nums2) {
		return Intersect(nums2, nums1)
	}
	m := map[string]int{}
	for _, num := range nums1 {
		m[num]++
	}
	var intersection []string
	for _, num := range nums2 {
		if m[num] > 0 {
			intersection = append(intersection, num)
			m[num]--
		}
	}
	return intersection
}
