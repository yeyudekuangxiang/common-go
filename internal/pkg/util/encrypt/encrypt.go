package encrypt

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"mio/config"
	"mio/internal/pkg/util"
	"net/url"
	"sort"
	"strings"
)

// GetSign 签名
func GetSign(params map[string]string) string {
	//排序
	var slice []string
	for k := range params {
		slice = append(slice, k)
	}
	sort.Strings(slice)
	var signStr string
	for _, v := range slice {
		signStr += PercentEncode(v) + "=" + PercentEncode(util.InterfaceToString(params[v])) + "&"
	}
	signStr = strings.TrimRight(signStr, "&")
	stringToSign := "POST&" + PercentEncode("/") + "&"
	stringToSign = stringToSign + PercentEncode(signStr)
	accessKeySecret := config.Config.ActivityZyh.AccessKeySecret + "&"
	return HMACSHA1(stringToSign, accessKeySecret)
}
func PercentEncode(string2 string) string {
	string2 = url.QueryEscape(string2)
	string2 = strings.Replace(string2, "+", "%20", -1)
	string2 = strings.Replace(string2, "*", "%2A", -1)
	string2 = strings.Replace(string2, "%7E", "~", -1)
	println(string2)
	return string2
}

/*
 keyStr 密钥
 value  消息内容
*/

func HMACSHA1(keyStr, value string) string {
	key := []byte(keyStr)
	mac := hmac.New(sha1.New, key)

	mac.Write([]byte(value))
	//进行base64编码
	res := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return res
}
