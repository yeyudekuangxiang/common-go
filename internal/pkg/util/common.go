package util

import (
	"github.com/shopspring/decimal"
	"math"
	"math/rand"
	"strings"
	"time"
)

func LinkJoin(ele ...string) string {
	builder := strings.Builder{}
	length := len(ele) - 1
	for i, e := range ele {

		if i == length {
			e = strings.TrimLeft(e, "/")
		} else {
			e = strings.Trim(e, "/")
		}
		builder.WriteString(e)

		if i != length {
			builder.WriteString("/")
		}
	}
	return builder.String()
}

func MapInterface2int64(inputData map[string]interface{}) map[string]int64 {
	outputData := map[string]int64{}
	for key, value := range inputData {
		switch value.(type) {
		case int64:
			outputData[key] = value.(int64)
		case string:
			outputData[key] = value.(int64)
		case int:
			outputData[key] = value.(int64)
		}
	}
	return outputData
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
	a1 := toRadians(decimal.NewFromFloat(l.Lat))
	b1 := toRadians(decimal.NewFromFloat(l.Lng))
	a2 := toRadians(decimal.NewFromFloat(point.Lat))
	b2 := toRadians(decimal.NewFromFloat(point.Lng))
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

func toRadians(num decimal.Decimal) decimal.Decimal {
	return num.Mul(decimal.NewFromFloat(math.Pi)).Div(decimal.NewFromInt32(180))
}

// Rand4Number 生成一个随机四位数
func Rand4Number() string {
	return string(rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
	//生成一个rand
}
