package signtool

import (
	"github.com/yeyudekuangxiang/common-go/tool/encrypttool"
	"sort"
	"strconv"
	"strings"
)

/**
签名规范(对外签名规范)
sign = md5({secret}param1={paramValue1}&param2={paramValue2}
1.
参数key升序排序 （ch={ch}&createdAt={createdAt}） 是案例，根据实际参数
2.secret 由绿喵提供,测试(0tlrEVZtRE)，正式会上生产再提供

举例：
参数： mobile：18800001111
参数整个成字符串 str = 0tlrEVZtREmobile=18800001111
对str进行md5加密，结果为：143D3F2E556798DEA6C6F34DB93F24B2
*/

func GetSign(params map[string]interface{}, key string, joiner string) string {
	var slice []string
	for k := range params {
		slice = append(slice, k)
	}
	sort.Strings(slice)
	var signStr string
	for _, v := range slice {
		signStr += v + "=" + InterfaceToString(params[v]) + joiner
	}
	signStr = strings.TrimRight(signStr, joiner)
	//验证签名
	return strings.ToUpper(encrypttool.Md5(key + signStr))
}

func InterfaceToString(data interface{}) string {
	var key string
	switch data.(type) {
	case string:
		key = data.(string)
	case int:
		key = strconv.Itoa(data.(int))
	case int64:
		it := data.(int64)
		key = strconv.FormatInt(it, 10)
	case float64:
		it := data.(float64)
		key = strconv.FormatFloat(it, 'f', -1, 64)
	}
	return key
}
