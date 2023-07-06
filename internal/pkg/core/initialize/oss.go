package initialize

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/aliyuntool"
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
)

func InitOss() {
	log.Println("初始化阿里云oss组件...")
	provider, err := aliyuntool.NewOssCredentialProvider(nil)
	if err != nil && err != aliyuntool.ErrCredentialNotFound {
		log.Panic("初始化阿里云ossProvider异常", err)
	}
	client, err := oss.New(config.Config.OSS.Endpoint, config.Config.OSS.AccessKey, config.Config.OSS.AccessSecret, oss.SetCredentialsProvider(provider))
	if err != nil {
		log.Panic(err)
	}
	*app.OssClient = *client
	log.Println("初始化阿里云oss组件成功")

}
func InitSts() {
	log.Println("初始化阿里云sts组件...")

	providerClient, err := sts.NewClientWithProvider(config.Config.OSS.Region)
	if err != nil && err.Error() != "No credential found" {
		log.Panic("初始化阿里云sts provider异常", err)
	}
	if providerClient != nil {
		*app.STSClient = *providerClient
	} else {
		client, err := sts.NewClientWithAccessKey(config.Config.OSS.Region, config.Config.OSS.AccessKey, config.Config.OSS.AccessSecret)
		if err != nil {
			log.Panic(err)
		}
		*app.STSClient = *client
	}

	log.Println("初始化阿里云sts组件成功")
}
