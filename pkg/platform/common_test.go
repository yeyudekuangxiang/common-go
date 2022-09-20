package platform

import (
	"fmt"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/encrypt"
	"sort"
	"testing"
)

// CheckSign 验证签名
func TestCheckSign(t *testing.T) {
	params := map[string]string{
		"platformKey": "zcyp",
		"mobile":      "15512341234",
		"method":      "apply",
		"sign":        "2a98658293181b71ca85d44b5b655249",
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
	//验证签名
	signMd5 := encrypt.Md5(params["platformKey"] + signStr)
	if signMd5 != sign {
		fmt.Printf("error --- oriSign: %s\nsign:%s\n", sign, signMd5)
	}
	fmt.Printf("oriSign: %s\nsign:%s\n", sign, signMd5)
}
