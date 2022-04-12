package duiba

import (
	mmd5 "crypto/md5"
	"fmt"
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
	return md5(signStr)
}
func md5(str string) string {
	return fmt.Sprintf("%x", mmd5.Sum([]byte(str)))
}
