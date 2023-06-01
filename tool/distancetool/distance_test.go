package distancetool

import (
	"fmt"
	"testing"
)

// 实现分页方法
func (rs resultSlice) Page(page int, pageSize int) []result {
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > len(rs) {
		end = len(rs)
	}
	if start > end {
		start = end
	}
	return rs[start:end]
}

func TestDis(t *testing.T) {
	// 起点
	origin := Point{lat: 31.23, lng: 121.47}

	// 目标点列表
	var targets []Point = []Point{
		{name: "北京", lat: 39.90, lng: 116.41},
		{name: "香港", lat: 22.27, lng: 114.17},
		{name: "广州", lat: 23.13, lng: 113.27},
		{name: "南京", lat: 32.05, lng: 118.78},
		{name: "上海", lat: 31.23, lng: 121.47},
		{name: "深圳", lat: 22.54, lng: 114.06},
		{name: "福州", lat: 26.08, lng: 119.30},
	}

	slice, err := DistanceArr(origin, targets)
	if err != nil {
		return
	}
	/*for _, r := range slice {
		fmt.Printf("%-4s距离 %.2f 米\n", r.point.name, r.distance)
	}*/

	// 分页
	pageSize := 2
	pageCount := (len(slice) + pageSize - 1) / pageSize
	for i := 1; i <= pageCount; i++ {
		pageResults := slice.Page(i, pageSize)
		fmt.Printf("第 %d 页：\n", i)
		for _, r := range pageResults {
			fmt.Printf("%-4s距离 %.2f 米\n", r.point.name, r.distance*1000)
		}
	}
	/*
		// 计算距离并存储结果到resultSlice
		var results resultSlice
		for _, target := range targets {
			d := distance(origin, target)
			r := result{distance: d, point: target}
			results = append(results, r)
		}

		// 按距离升序排序
		sort.Sort(results)

		// 输出结果
		for _, r := range results {
			fmt.Printf("%-4s距离 %.2f 米\n", r.point.name, r.distance)
		}*/

}
