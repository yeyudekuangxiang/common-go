package service

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"hash"
	"io"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/service/service_types"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"strings"
	"time"
)

const OssDomain = "https://miotech-resource.oss-cn-hongkong.aliyuncs.com"

var DefaultOssService OssService

func InitDefaultOssService() {
	DefaultOssService = OssService{
		client:   app.OssClient,
		Bucket:   "miotech-resource",
		BasePath: config.Config.OSS.BasePath,
	}
}

type OssService struct {
	client   *oss.Client
	Bucket   string
	BasePath string
}

// PutObject name 文件路径  最终路径为 OssService.BasePath +"/"+ images/topic/tag/1.png
func (srv OssService) PutObject(name string, reader io.Reader) (string, error) {
	name = srv.BasePath + "/" + strings.TrimLeft(name, "/")
	return srv.PubObjectAbsolutePath(name, reader)
}

// PubObjectAbsolutePath name 文件路径 例如static/mp2c/images/topic/tag/1.png
func (srv OssService) PubObjectAbsolutePath(name string, reader io.Reader) (string, error) {
	name = strings.TrimLeft(name, "/")
	bucket, err := srv.client.Bucket(srv.Bucket)
	if err != nil {
		return "", err
	}
	err = bucket.PutObject(name, reader)
	if err != nil {
		return "", err
	}
	return util.LinkJoin(OssDomain, name), nil
}
func (srv OssService) GetPolicyToken(param service_types.GetOssPolicyTokenParam) (*service_types.OssPolicyToken, error) {
	expireEnd := time.Now().Add(param.ExpireTime).Unix()
	tokenExpire := time.Unix(expireEnd, 0).UTC().Format("2006-01-02T15:04:05Z")

	//create post policy json
	policyConfig := service_types.OssPolicyConfig{}
	policyConfig.Expiration = tokenExpire
	policyConfig.AddStartWithKey(param.UploadDir)
	policyConfig.AddContentLength(param.MaxSize)
	policyConfig.AddBucket(srv.Bucket)
	policyConfig.AddContentType(param.MimeTypes)

	//calucate signature
	configData, err := json.Marshal(policyConfig)
	if err != nil {
		return nil, errno.ErrInternalServer.With(err)
	}

	baseConfig := base64.StdEncoding.EncodeToString(configData)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(config.Config.OSS.AccessSecret))

	_, err = io.WriteString(h, baseConfig)
	if err != nil {
		return nil, errno.ErrInternalServer.With(err)
	}

	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	callbackParam := service_types.OssCallbackParam{}
	callbackParam.CallbackUrl = param.CallbackUrl
	callbackParam.CallbackBody = "filename=${object}&size=${size}&mimeType=${mimeType}&height=${imageInfo.height}&width=${imageInfo.width}"
	callbackParam.CallbackBodyType = "application/x-www-form-urlencoded"

	callbackStr, err := json.Marshal(callbackParam)
	if err != nil {
		return nil, errno.ErrInternalServer.With(err)
	}
	callbackBase64 := base64.StdEncoding.EncodeToString(callbackStr)

	policyToken := service_types.OssPolicyToken{}
	policyToken.AccessKeyId = config.Config.OSS.AccessKey
	policyToken.Host = "https://" + srv.Bucket + "." + config.Config.OSS.Endpoint
	policyToken.Expire = expireEnd
	policyToken.Signature = signedStr
	policyToken.Directory = param.UploadDir
	policyToken.Policy = baseConfig
	policyToken.Callback = callbackBase64

	return &policyToken, nil
}
