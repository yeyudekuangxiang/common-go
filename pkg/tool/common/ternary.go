package common

import "reflect"

type TOResult struct {
	val interface{}
}

func (res TOResult) String() string {
	return res.val.(string)
}
func (res TOResult) Int() int {
	return int(res.Int64())
}
func (res TOResult) Int64() int64 {
	return reflect.ValueOf(res.val).Int()
}
func (res TOResult) Float32() float32 {
	return float32(res.Float64())
}
func (res TOResult) Float64() float64 {
	return reflect.ValueOf(res.val).Float()
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
