package distancetool

import (
	"fmt"
	"github.com/shopspring/decimal"
	"math"
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
	a := CalcLngLatDistance(116.4133836971231, 39.910924547299565, 121.48053886017651, 31.235929042252014)

	println(a)

	b := Distance(Point{
		Lat: 39.910924547299565,
		Lng: 116.4133836971231,
	}, Point{
		Lat: 31.235929042252014,
		Lng: 121.48053886017651,
	})
	println(b)
	// 起点
	origin := PointArr{Lat: 31.23, Lng: 121.47}

	// 目标点列表
	var targets []PointArr = []PointArr{
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

// CalcLngLatDistance 根据两点的经纬度计算直线距离 返回米 保留两位小数
// 计算北京到上海距离 CalcLngLatDistance(116.4133836971231, 39.910924547299565, 121.48053886017651,31.235929042252014)
func CalcLngLatDistance(lng1 float64, lat1 float64, lng2 float64, lat2 float64) float64 {
	return LatLon{Lat: lat1, Lng: lng1}.DistanceTo(LatLon{Lat: lat2, Lng: lng2})
}

type LatLon struct {
	Lng float64
	Lat float64
}

// DistanceTo 根据两点的经纬度计算直线距离 返回米 保留两位小数
// 计算北京到上海距离 LatLon{Lat: 39.910924547299565, Lng: 116.4133836971231}.DistanceTo(LatLon{Lat: 31.235929042252014, Lng: 121.48053886017651})
// 算法参考 https://github.com/chrisveness/geodesy
func (l LatLon) DistanceTo(point LatLon) float64 {
	R := 6371e3
	a1 := toRadiansV2(decimal.NewFromFloat(l.Lat))
	b1 := toRadiansV2(decimal.NewFromFloat(l.Lng))
	a2 := toRadiansV2(decimal.NewFromFloat(point.Lat))
	b2 := toRadiansV2(decimal.NewFromFloat(point.Lng))
	ca := a2.Sub(a1)
	cb := b2.Sub(b1)

	s1 := ca.Div(decimal.NewFromInt(2)).Sin().Pow(decimal.NewFromInt(2))
	s2 := a1.Cos().Mul(a2.Cos())
	s3 := cb.Div(decimal.NewFromInt(2)).Sin().Pow(decimal.NewFromInt(2))

	a := s1.Add(s2.Mul(s3))

	c := decimal.NewFromFloat(math.Atan2(math.Sqrt(a.InexactFloat64()), math.Sqrt(decimal.NewFromInt(1).Sub(a).InexactFloat64()))).Mul(decimal.NewFromInt32(2))
	d := c.Mul(decimal.NewFromFloat(R))
	return d.Round(2).InexactFloat64()
}

func toRadiansV2(num decimal.Decimal) decimal.Decimal {
	return num.Mul(decimal.NewFromFloat(math.Pi)).Div(decimal.NewFromInt32(180))
}
