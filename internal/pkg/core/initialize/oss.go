package initialize

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
)

func initOss() {
	client, err := oss.New(config.Config.OSS.Endpoint, config.Config.OSS.AccessKey, config.Config.OSS.AccessSecret)
	if err != nil {
		log.Fatal(err)
	}
	*app.OssClient = *client
}
