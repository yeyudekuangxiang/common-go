package util

import (
	"strings"
)

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
func IntersectContains(identify []string, rules []string) map[string]string {
	m := map[string]string{}
	for i, v := range identify {
		for _, rule := range rules {
			if strings.Contains(v, rule) {
				//匹配字符 到map
				val := strings.Trim(v, rule)
				if val == "" {
					val = identify[i+1]
				}
				m[rule] = val
			}
		}
	}
	return m
}

// SortArray 快速排序，基于比较，不稳定算法
func SortArray(nums []int) []int {
	var quick func(nums []int, left, right int) []int
	quick = func(nums []int, left, right int) []int {
		// 递归终止条件
		if left > right {
			return nil
		}
		// 左右指针及主元
		i, j, pivot := left, right, nums[left]
		for i < j {
			// 寻找小于主元的右边元素
			for i < j && nums[j] >= pivot {
				j--
			}
			// 寻找大于主元的左边元素
			for i < j && nums[i] <= pivot {
				i++
			}
			// 交换i/j下标元素
			nums[i], nums[j] = nums[j], nums[i]
		}
		// 交换元素
		nums[i], nums[left] = nums[left], nums[i]
		quick(nums, left, i-1)
		quick(nums, i+1, right)
		return nums
	}
	return quick(nums, 0, len(nums)-1)
}
