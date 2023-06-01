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
	origin := Point{Lat: 31.23, Lng: 121.47}

	// 目标点列表
	var targets []Point = []Point{
		{Name: "北京", Lat: 39.90, Lng: 116.41},
		{Name: "香港", Lat: 22.27, Lng: 114.17},
		{Name: "广州", Lat: 23.13, Lng: 113.27},
		{Name: "南京", Lat: 32.05, Lng: 118.78},
		{Name: "上海", Lat: 31.23, Lng: 121.47},
		{Name: "深圳", Lat: 22.54, Lng: 114.06},
		{Name: "福州", Lat: 26.08, Lng: 119.30},
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
			fmt.Printf("%-4s距离 %.2f 米\n", r.Point.Name, r.Distance*1000)
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
