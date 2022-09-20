package platform

import (
	"errors"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/encrypt"
	"sort"
)

// CheckSign 验证签名
func CheckSign(params map[string]interface{}) error {
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
	//验证签名
	signMd5 := encrypt.Md5(params["platformKey"].(string) + signStr)
	if signMd5 != sign {
		return errors.New("验签失败 oriSign: " + sign.(string) + " encodeSign: " + signMd5)
	}
	return nil
}
