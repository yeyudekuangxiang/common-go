package util

import (
	mmd5 "crypto/md5"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/sorttool"
	"net/url"
)

func MapTo(data interface{}, v interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  v,
		TagName: "json",
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(decoder.Decode(data))
}
func BuildQuery(params map[string]string) string {
	vals := url.Values{}
	sorttool.Map(params, func(key interface{}) {
		k := key.(string)
		if k == "" {
			return
		}
		vals.Set(k, params[k])
	})
	return vals.Encode()
}
func Assign(m1 map[string]string, m2 map[string]string) map[string]string {
	m := make(map[string]string)
	for k, v := range m1 {
		m[k] = v
	}
	for k, v := range m2 {
		m[k] = v
	}
	return m
}
func Md5(str string) string {
	return fmt.Sprintf("%x", mmd5.Sum([]byte(str)))
}
