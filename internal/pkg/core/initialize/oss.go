package initialize

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
)

func InitOss() {
	log.Println("初始化阿里云oss组件...")
	client, err := oss.New(config.Config.OSS.Endpoint, config.Config.OSS.AccessKey, config.Config.OSS.AccessSecret)
	if err != nil {
		log.Panic(err)
	}
	*app.OssClient = *client
	log.Println("初始化阿里云oss组件成功")

}
func InitSts() {
	log.Println("初始化阿里云sts组件...")
	client, err := sts.NewClientWithAccessKey("", config.Config.OSS.AccessKey, config.Config.OSS.AccessSecret)
	if err != nil {
		log.Panic(err)
	}
	*app.STSClient = *client
	log.Println("初始化阿里云sts组件成功")

}
