package distancetool

import (
	"github.com/pkg/errors"
	"math"
	"sort"
)

const earthRadius = 6371000 // 地球赤道半径，单位米

// 点的结构体，包含经纬度和点的名称

type Point struct {
	Lat float64
	Lng float64
}

type PointArr struct {
	Name string
	Lat  float64
	Lng  float64
}

// 用于计算距离的函数

func Distance(p1, p2 Point) float64 {
	deltaLat := toRadians(p2.Lat - p1.Lat)
	deltaLng := toRadians(p2.Lng - p1.Lng)
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(toRadians(p1.Lat))*math.Cos(toRadians(p2.Lat))*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}

// 将角度转换为弧度
func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func DistanceArr(origin PointArr, targets []PointArr) (resultSlice, error) {
	if len(targets) == 0 {
		return nil, errors.New("参数有误")
	}
	// 计算距离并存储结果到resultSlice
	var results resultSlice
	for _, target := range targets {
		d := Distance(Point{
			Lat: origin.Lat,
			Lng: origin.Lng,
		}, Point{
			Lat: target.Lat,
			Lng: target.Lng,
		})
		r := result{Distance: d, Point: target}
		results = append(results, r)
	}
	// 按距离升序排序
	sort.Sort(results)
	return results, nil
}

// 结果项的结构体，包含距离和点的信息
type result struct {
	Distance float64
	Point    PointArr
}

// 定义结果项列表，用于存储每个点与起点之间的距离
type resultSlice []result

// 实现sort.Interface接口的Len方法

func (rs resultSlice) Len() int {
	return len(rs)
}

// 实现sort.Interface接口的Swap方法

func (rs resultSlice) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

// 实现sort.Interface接口的Less方法，按距离升序排列

func (rs resultSlice) Less(i, j int) bool {
	return rs[i].Distance < rs[j].Distance
}
