package distancetool

import (
	"fmt"
	"github.com/pkg/errors"
	"math"
	"sort"
)

const earthRadius = 6371000 // 地球赤道半径，单位米

// 点的结构体，包含经纬度

type Point struct {
	Lat float64
	Lng float64
}

//点的结构体，包含经纬度和点的名称

type PointArr struct {
	Name string
	Lat  float64
	Lng  float64
}

// 将角度转换为弧度

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// 计算距离的函数

func Distance(p1, p2 Point) float64 {
	deltaLat := toRadians(p2.Lat - p1.Lat)
	deltaLng := toRadians(p2.Lng - p1.Lng)
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(toRadians(p1.Lat))*math.Cos(toRadians(p2.Lat))*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}

// 计算一个点到多个点距离列表

func DistanceArr(origin PointArr, targets []PointArr) (ResultSlice, error) {
	if len(targets) == 0 {
		return nil, errors.New("参数有误")
	}
	// 计算距离并存储结果到resultSlice
	var results ResultSlice
	for _, target := range targets {
		d := Distance(Point{
			Lat: origin.Lat,
			Lng: origin.Lng,
		}, Point{
			Lat: target.Lat,
			Lng: target.Lng,
		})
		r := Result{Distance: d, Point: target}
		results = append(results, r)
	}
	return results, nil
}

//距离单位转换

func FormatDistance(distance float64) string {
	if distance < 1000 {
		return fmt.Sprintf("%.2f m", distance)
	} else {
		return fmt.Sprintf("%.2f km", distance/1000.0)
	}
}

// 结果项的结构体，包含距离和点的信息

type Result struct {
	Distance float64
	Point    PointArr
}

// 定义结果项列表，用于存储每个点与起点之间的距离

type ResultSlice []Result

// 实现sort.Interface接口的Len方法

func (rs ResultSlice) Len() int {
	return len(rs)
}

// 实现sort.Interface接口的Swap方法

func (rs ResultSlice) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

// 实现sort.Interface接口的Less方法，按距离升序排列

func (rs ResultSlice) Less(i, j int) bool {
	return rs[i].Distance < rs[j].Distance
}

//距离列表，按距离升序排序

func (rs ResultSlice) DistanceSort() (ResultSlice, error) {
	sort.Sort(rs)
	return rs, nil
}

// 实现分页方法

func (rs ResultSlice) Page(page int, pageSize int) []Result {
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
