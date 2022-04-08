package util

import (
	"github.com/pkg/errors"
	"mio/internal/pkg/core/app"
	"os"
	"strings"
)

func GetAppConfig(key string) (string, error) {
	if key == "" {
		return "", errors.New("key can not be empty")
	}
	keys := strings.Split(key, ".")
	if len(keys) == 1 {
		return "", errors.New("need section")
	}
	section := app.Ini.Section(keys[0])
	if section == nil {
		return "", errors.New("not found section")
	}
	item := section.Key(keys[1])
	if item == nil {
		return "", errors.New("not found key")
	}
	return item.Value(), nil
}
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

// TernaryOperator 三元运算
func TernaryOperator(right bool, val1 interface{}, val2 interface{}) TOResult {
	if right {
		return TOResult{val: val1}
	}
	return TOResult{val: val2}
}
