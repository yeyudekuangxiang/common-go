package util

import (
	"os"
	"strings"
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
