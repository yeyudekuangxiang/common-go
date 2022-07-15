package service_types

import (
	"time"
)

type GetOssPolicyTokenParam struct {
	//有效期
	ExpireTime time.Duration
	//单位 B
	MaxSize     int64
	UploadDir   string
	CallbackUrl string
	MimeTypes   []string
	MaxAge      int //缓存有效期 单位秒
}
type OssPolicyConfig struct {
	Expiration string        `json:"expiration"`
	Conditions []interface{} `json:"conditions"`
}

//AddContentLength 上传文件大小 单位 B 1GB=1024MB，1MB=1024KB，1KB=1024B
func (c *OssPolicyConfig) AddContentLength(maxLength int64) {
	c.Conditions = append(c.Conditions, []interface{}{
		"content-length-range",
		1,
		maxLength,
	})
}

// AddBucket 添加bucket限制
func (c *OssPolicyConfig) AddBucket(bucket string) {
	c.Conditions = append(c.Conditions, map[string]string{
		"bucket": bucket,
	})
}

func (c *OssPolicyConfig) AddMaxAge(maxAge int) {
	/*c.Conditions = append(c.Conditions, []interface{}{
		"not-in",
		"Cache-Control",
		[]string{"no-cache"},
	})*/
}

//AddContentType []string{"image/jpg","image/png"}
func (c *OssPolicyConfig) AddContentType(contentTypes []string) {
	if len(contentTypes) == 0 {
		return
	}
	c.Conditions = append(c.Conditions, []interface{}{
		"in",
		"$content-type",
		contentTypes,
	})
}

//AddStartWithKey 上传文件目录 mp2c/user/avatar
func (c *OssPolicyConfig) AddStartWithKey(dir string) {
	c.Conditions = append(c.Conditions, []interface{}{
		"starts-with",
		"$key",
		dir,
	})
}

type OssCallbackParam struct {
	CallbackUrl      string `json:"callbackUrl"`
	CallbackBody     string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}

type OssPolicyToken struct {
	AccessKeyId string `json:"AccessKeyId"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Directory   string `json:"dir"`
	Callback    string `json:"callback"`
}
