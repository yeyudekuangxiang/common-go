package platform

import (
	"github.com/pkg/errors"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/encrypt"
	"sort"
	"strings"
)

// CheckSign 验证签名
func CheckSign(sign string, params map[string]interface{}, key string, joiner string) error {
	signMd5 := GetSign(params, key, joiner)
	//验证签名
	if signMd5 != sign {
		return errors.New("验签失败")
	}
	return nil
}

// GetSign 签名
func GetSign(params map[string]interface{}, key string, joiner string) string {
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
	return encrypt.Md5(key + signStr)
}
