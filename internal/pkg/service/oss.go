package service

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"mio/config"
	"mio/internal/pkg/core/app"
	"strings"
)

const OssDomain = "https://miotech-resource.oss-cn-hongkong.aliyuncs.com/"

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
	return OssDomain + name, nil
}
