package aliyuntool

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"hash"
	"io"
	"time"
)

type OssClient struct {
	Client *oss.Client
}

func NewOssClient(client *oss.Client) *OssClient {
	return &OssClient{Client: client}
}

type CreateOssPolicyTokenParam struct {
	Bucket string
	//有效期
	ExpireTime time.Duration
	//单位 B >=1
	MinSize int64
	//单位 B
	MaxSize     int64
	UploadDir   string
	CallbackUrl string
	MimeTypes   []string
	MaxAge      int64 //缓存有效期 单位秒
}
type CreatePolicyTokenResult struct {
	AccessKeyId       string
	Host              string
	Expire            int64
	Signature         string
	Policy            string
	Directory         string
	Callback          string
	XOssSecurityToken string
}

type OssPolicyConfig struct {
	Expiration string        `json:"expiration"`
	Conditions []interface{} `json:"conditions"`
}

// AddContentLength 上传文件大小 单位 B 1GB=1024MB，1MB=1024KB，1KB=1024B
func (c *OssPolicyConfig) AddContentLength(minLength, maxLength int64) {
	c.Conditions = append(c.Conditions, []interface{}{
		"content-length-range",
		minLength,
		maxLength,
	})
}

// AddBucket 添加bucket限制
func (c *OssPolicyConfig) AddBucket(bucket string) {
	c.Conditions = append(c.Conditions, map[string]string{
		"bucket": bucket,
	})
}

func (c *OssPolicyConfig) AddMaxAge(maxAge int64) {
	/*c.Conditions = append(c.Conditions, []interface{}{
		"in",
		"$cache-control",
		[]string{"max-age=" + strconv.FormatInt(maxAge, 10)},
	})*/
}

// AddContentType []string{"image/jpg","image/png"}
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

// AddStartWithKey 上传文件目录 mp2c/user/avatar
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

func (c *OssClient) CreatePolicyToken(param CreateOssPolicyTokenParam) (*CreatePolicyTokenResult, error) {
	if param.MinSize <= 0 {
		param.MinSize = 1
	}
	expireEnd := time.Now().Add(param.ExpireTime).Unix()
	tokenExpire := time.Unix(expireEnd, 0).UTC().Format("2006-01-02T15:04:05Z")

	//create post policy json
	policyConfig := OssPolicyConfig{}
	policyConfig.Expiration = tokenExpire
	policyConfig.AddStartWithKey(param.UploadDir)
	policyConfig.AddContentLength(param.MinSize, param.MaxSize)
	policyConfig.AddBucket(param.Bucket)
	policyConfig.AddMaxAge(param.MaxAge)
	policyConfig.AddContentType(param.MimeTypes)

	//calucate signature
	configData, err := json.Marshal(policyConfig)
	if err != nil {
		return nil, err
	}

	accessSecret := c.Client.Config.GetCredentials().GetAccessKeySecret()

	baseConfig := base64.StdEncoding.EncodeToString(configData)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(accessSecret))

	_, err = io.WriteString(h, baseConfig)
	if err != nil {
		return nil, err
	}

	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	callbackParam := OssCallbackParam{}
	callbackParam.CallbackUrl = param.CallbackUrl
	callbackParam.CallbackBody = "filename=${object}&size=${size}&mimeType=${mimeType}&height=${imageInfo.height}&width=${imageInfo.width}"
	callbackParam.CallbackBodyType = "application/x-www-form-urlencoded"

	callbackStr, err := json.Marshal(callbackParam)
	if err != nil {
		return nil, err
	}
	callbackBase64 := base64.StdEncoding.EncodeToString(callbackStr)

	policyToken := CreatePolicyTokenResult{}
	policyToken.AccessKeyId = c.Client.Config.GetCredentials().GetAccessKeyID()
	policyToken.Host = "https://" + param.Bucket + "." + c.Client.Config.Endpoint
	policyToken.Expire = expireEnd
	policyToken.Signature = signedStr
	policyToken.Directory = param.UploadDir
	policyToken.Policy = baseConfig
	policyToken.Callback = callbackBase64
	policyToken.XOssSecurityToken = c.Client.Config.GetCredentials().GetSecurityToken()

	return &policyToken, nil
}
