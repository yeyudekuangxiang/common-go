package initialize

import (
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
