package platform

import (
	"errors"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/encrypt"
	"sort"
	"strings"
)

// CheckSign 验证签名
func CheckSign(params map[string]interface{}, joiner string) error {
	if joiner == "" {
		joiner = ";"
	}
	sign := params["sign"]
	delete(params, "sign")
	var slice []string
	for k := range params {
		slice = append(slice, k)
	}
	sort.Strings(slice)
	var signStr string
	for _, v := range slice {
		signStr += v + "=" + util.InterfaceToString(params[v]) + ";"
	}
	if joiner != ";" {
		strings.TrimRight(signStr, joiner)
	}
	//验证签名
	signMd5 := encrypt.Md5(params["platformKey"].(string) + signStr)
	if signMd5 != sign {
		app.Logger.Errorf("验签失败 oriSign: %s ; encodeSign: %s", sign, signMd5)
		return errors.New("验签失败")
	}
	return nil
}

// GetSign 签名
func GetSign(params map[string]interface{}, joiner string) string {
	var platformKey string
	if key, ok := params["platformKey"]; ok {
		platformKey = key.(string)
	}
	if joiner == "" {
		joiner = ";"
	}
	var slice []string
	for k := range params {
		slice = append(slice, k)
	}
	sort.Strings(slice)
	var signStr string
	for _, v := range slice {
		signStr += v + "=" + util.InterfaceToString(params[v]) + joiner
	}
	if joiner != ";" {
		signStr = strings.TrimRight(signStr, joiner)
	}
	//验证签名
	return encrypt.Md5(platformKey + signStr)
}
