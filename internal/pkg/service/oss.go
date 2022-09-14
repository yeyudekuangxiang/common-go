package service

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"hash"
	"io"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"os"
	"strings"
	"time"
)

var DefaultOssService OssService

func InitDefaultOssService() {
	DefaultOssService = OssService{
		client:   app.OssClient,
		Bucket:   "miotech-resource",
		BasePath: config.Config.OSS.BasePath,
		Domain:   config.Config.OSS.CdnDomain,
	}
}

type OssService struct {
	client   *oss.Client
	Bucket   string
	BasePath string
	Domain   string
}

// PutObject name 文件路径  最终路径为 OssService.BasePath +"/"+ images/topic/tag/1.png 上传成功返回路径
func (srv OssService) PutObject(name string, reader io.Reader) (string, error) {
	name = srv.BasePath + "/" + strings.TrimLeft(name, "/")
	return srv.PubObjectAbsolutePath(name, reader)
}

// PubObjectAbsolutePath name 文件路径 例如static/mp2c/images/topic/tag/1.png 上传成功返回路径
func (srv OssService) PubObjectAbsolutePath(name string, reader io.Reader) (string, error) {
	start := time.Now()
	defer func() {
		app.Logger.Infof("上传文件耗时 %s %s", name, time.Now().Sub(start).String())
	}()
	name = strings.TrimLeft(name, "/")
	bucket, err := srv.client.Bucket(srv.Bucket)
	if err != nil {
		return "", err
	}
	err = bucket.PutObject(name, reader)
	if err != nil {
		return "", err
	}
	return name, nil
}
func (srv OssService) GetPolicyToken(param srv_types.GetOssPolicyTokenParam) (*srv_types.OssPolicyToken, error) {
	expireEnd := time.Now().Add(param.ExpireTime).Unix()
	tokenExpire := time.Unix(expireEnd, 0).UTC().Format("2006-01-02T15:04:05Z")

	//create post policy json
	policyConfig := srv_types.OssPolicyConfig{}
	policyConfig.Expiration = tokenExpire
	policyConfig.AddStartWithKey(param.UploadDir)
	policyConfig.AddContentLength(param.MaxSize)
	policyConfig.AddBucket(srv.Bucket)
	policyConfig.AddMaxAge(param.MaxAge)
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

	callbackParam := srv_types.OssCallbackParam{}
	callbackParam.CallbackUrl = param.CallbackUrl
	callbackParam.CallbackBody = "filename=${object}&size=${size}&mimeType=${mimeType}&height=${imageInfo.height}&width=${imageInfo.width}"
	callbackParam.CallbackBodyType = "application/x-www-form-urlencoded"

	callbackStr, err := json.Marshal(callbackParam)
	if err != nil {
		return nil, errno.ErrInternalServer.With(err)
	}
	callbackBase64 := base64.StdEncoding.EncodeToString(callbackStr)

	policyToken := srv_types.OssPolicyToken{}
	policyToken.AccessKeyId = config.Config.OSS.AccessKey
	policyToken.Host = "https://" + srv.Bucket + "." + config.Config.OSS.Endpoint
	policyToken.Expire = expireEnd
	policyToken.Signature = signedStr
	policyToken.Directory = param.UploadDir
	policyToken.Policy = baseConfig
	policyToken.Callback = callbackBase64

	return &policyToken, nil
}
func (srv OssService) FullUrl(path string) string {
	return util.LinkJoin(srv.Domain, path)
}

func (srv OssService) MultipartPutObject(name string, reader io.Reader, locaFilename string) (string, error) {
	// 创建OSSClient实例。
	// yourEndpoint填写Bucket所在地域对应的Endpoint。以华东1（杭州）为例，Endpoint填写为https://oss-cn-hangzhou.aliyuncs.com。
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	// 填写存储空间名称。
	bucketName := srv.Bucket
	// 填写Object完整路径。Object完整路径中不能包含Bucket名称。
	objectName := name
	// 填写本地文件的完整路径。如果未指定本地路径，则默认从示例程序所属项目对应本地路径中上传文件。
	locaFilename = "/Users/apple/Desktop/lm1.mp4"
	// 获取存储空间。
	bucket, err := srv.client.Bucket(bucketName)
	if err != nil {
		return "", err
		os.Exit(-1)
	}
	// 将本地文件分片，且分片数量指定为3。
	chunks, err := oss.SplitFileByPartNum(locaFilename, 3)
	fd, err := os.Open(locaFilename)
	defer fd.Close()
	// 指定过期时间。
	expires := time.Date(2049, time.January, 10, 23, 0, 0, 0, time.UTC)
	// 如果需要在初始化分片时设置请求头，请参考以下示例代码。
	options := []oss.Option{
		oss.MetadataDirective(oss.MetaReplace),
		oss.Expires(expires),
		// 指定该Object被下载时的网页缓存行为。
		// oss.CacheControl("no-cache"),
		// 指定该Object被下载时的名称。
		// oss.ContentDisposition("attachment;filename=FileName.txt"),
		// 指定该Object的内容编码格式。
		// oss.ContentEncoding("gzip"),
		// 指定对返回的Key进行编码，目前支持URL编码。
		// oss.EncodingType("url"),
		// 指定Object的存储类型。
		// oss.ObjectStorageClass(oss.StorageStandard),
	}

	// 步骤1：初始化一个分片上传事件，并指定存储类型为标准存储。
	imur, err := bucket.InitiateMultipartUpload(objectName, options...)
	// 步骤2：上传分片。
	var parts []oss.UploadPart
	for _, chunk := range chunks {
		fd.Seek(chunk.Offset, os.SEEK_SET)
		// 调用UploadPart方法上传每个分片。
		part, err := bucket.UploadPart(imur, reader, chunk.Size, chunk.Number)
		if err != nil {
			return "", err
			os.Exit(-1)
		}
		parts = append(parts, part)
	}

	// 指定Object的读写权限为公共读，默认为继承Bucket的读写权限。
	objectAcl := oss.ObjectACL(oss.ACLPublicRead)

	// 步骤3：完成分片上传，指定文件读写权限为公共读。
	cmur, err := bucket.CompleteMultipartUpload(imur, parts, objectAcl)
	if err != nil {
		return "", err
		os.Exit(-1)
	}
	fmt.Println("cmur:", cmur)
	return cmur.Location, nil
}
