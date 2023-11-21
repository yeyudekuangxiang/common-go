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
	options := make([]oss.ClientOption, 0)
	if config.Config.OSS.AccessKey == "" || config.Config.OSS.AccessSecret == "" {
		provider, err := aliyuntool.NewOssCredentialProvider(nil)
		if err != nil && err != aliyuntool.ErrCredentialNotFound {
			log.Panic("初始化阿里云ossProvider异常", err)
		}
		if provider != nil {
			log.Printf("配置:%+v provider:%+v\n", config.Config.OSS, provider)
			log.Println("provider key ", provider.GetCredentials().GetAccessKeyID(), provider.GetCredentials().GetAccessKeySecret(), provider.GetCredentials().GetSecurityToken())
			options = append(options, oss.SetCredentialsProvider(provider))
		}
	}

	client, err := oss.New(config.Config.OSS.Endpoint, config.Config.OSS.AccessKey, config.Config.OSS.AccessSecret, options...)
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

	// 屏蔽wukey
	providerClient = nil

	if providerClient != nil {
		log.Printf("通过Provider初始化sts client %+v\n", providerClient)
		*app.STSClient = *providerClient
	} else {
		log.Println("通过accessKey初始化sts")
		client, err := sts.NewClientWithAccessKey(config.Config.OSS.Region, config.Config.OSS.AccessKey, config.Config.OSS.AccessSecret)
		if err != nil {
			log.Panic(err)
		}
		*app.STSClient = *client
	}

	log.Println("初始化阿里云sts组件成功")
}
