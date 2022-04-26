package service

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"mio/internal/pkg/core/app"
	"strings"
)

const OssDomain = "https://miotech-resource.oss-cn-hongkong.aliyuncs.com/"

var DefaultOssService = OssService{
	client: app.OssClient,
	Bucket: "miotech-resource",
}

type OssService struct {
	client *oss.Client
	Bucket string
}

// PutObject name 文件路径 static/mp2c/images/topic/tag/1.png
func (srv OssService) PutObject(name string, reader io.Reader) (string, error) {
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
