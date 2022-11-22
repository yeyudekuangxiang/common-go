package duiba

import (
	"gitlab.miotech.com/miotech-application/backend/common-go/pkg/duiba/util"
	"sort"
)

func sign(params map[string]string) string {
	keyList := make([]string, 0)
	for k := range params {
		keyList = append(keyList, k)
	}
	sort.Strings(keyList)

	signStr := ""
	for _, k := range keyList {
		if params[k] == "" {
			continue
		}
		signStr += params[k]
	}
	return util.Md5(signStr)
}
