package util

import (
	"github.com/shopspring/decimal"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

func IsTesting() bool {
	return os.Getenv("TEST_ENV") != ""
}

type TOResult struct {
	val interface{}
}

func (res TOResult) String() string {
	return res.val.(string)
}
func (res TOResult) Int() int {
	return res.val.(int)
}
func (res TOResult) Int64() int64 {
	return res.val.(int64)
}
func (res TOResult) Float32() float32 {
	return res.val.(float32)
}
func (res TOResult) Float64() float64 {
	return res.val.(float64)
}
func (res TOResult) Bool() bool {
	return res.val.(bool)
}
func (res TOResult) Interface() interface{} {
	return res.val
}

// Ternary 三元运算
func Ternary(right bool, val1 interface{}, val2 interface{}) TOResult {
	if right {
		return TOResult{val: val1}
	}
	return TOResult{val: val2}
}

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

// CalcLngLatDistance 根据经纬度计算距离 返回m
func CalcLngLatDistance(lng1 float64, lat1 float64, lng2 float64, lat2 float64) float64 {

	dlng1 := decimal.NewFromFloat(lng1)
	dlat1 := decimal.NewFromFloat(lat1)
	dlng2 := decimal.NewFromFloat(lng2)
	dlat2 := decimal.NewFromFloat(lat2)
	a := dlat1.Sub(dlat2).Div(decimal.NewFromInt32(2)).Sin().Pow(decimal.NewFromInt32(2))
	b := dlng1.Sub(dlng2).Div(decimal.NewFromInt32(2)).Sin().Pow(decimal.NewFromInt32(2))
	c := dlat1.Cos().Mul(dlat2.Cos()).Mul(b)
	f, _ := a.Add(c).Float64()

	distance := decimal.NewFromFloat(math.Asin(math.Sqrt(f))).Mul(decimal.NewFromInt32(2)).Mul(decimal.NewFromInt32(6378137))
	v, _ := distance.Round(2).Float64()
	return v
}

// Rand4Number 生成一个随机四位数
func Rand4Number() string {
	return string(rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
	//生成一个rand
}
