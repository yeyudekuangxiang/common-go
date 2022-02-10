package acm

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type Config struct {
	Endpoint    string
	NamespaceId string
	AccessKey   string
	SecretKey   string
	LogDir      string //the directory for log, default is current path
	RotateTime  string //the rotate time for log, eg: 30m, 1h, 24h, default is 24h
	MaxAge      int64  //the max age of a log file, default value is 3
	LogLevel    string //the level of log, it's must be debug,info,warn,error, default value is info
}

func NewClient(config Config) (*Client, error) {
	configClient, err := newConfigClient(config)
	if err != nil {
		return nil, err
	}
	return &Client{
		configClient,
	}, nil
}
func newConfigClient(config Config) (config_client.IConfigClient, error) {
	clientConfig := constant.ClientConfig{
		Endpoint:    config.Endpoint + ":8080",
		NamespaceId: config.NamespaceId,
		AccessKey:   config.AccessKey,
		SecretKey:   config.SecretKey,
		TimeoutMs:   5 * 1000,
		LogDir:      config.LogDir,
		RotateTime:  config.RotateTime,
		MaxAge:      config.MaxAge,
		LogLevel:    config.LogLevel,
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"clientConfig": clientConfig,
	})
	if err != nil {
		return nil, err
	}
	return configClient, nil
}

type Client struct {
	client config_client.IConfigClient
}

func (c Client) GetConfig(group string, dataId string) (content string, err error) {
	// 获取配置
	return c.client.GetConfig(vo.ConfigParam{
		Group:  group,
		DataId: dataId,
	})
}
