package config

import (
	"mio/internal/pkg/model/entity"
)

var Config = app{
	App:              appSetting{},
	Http:             httpSetting{},
	Database:         databaseSetting{},
	Log:              logSetting{},
	Weapp:            wxSetting{},
	MioSubOA:         wxSetting{},
	MioSrvOA:         wxSetting{},
	Redis:            redisSetting{},
	DuiBa:            duiBaSetting{},
	OSS:              ossSetting{},
	BaiDu:            baiDuSetting{},
	BaiDuImageSearch: baiDuImageSearchSetting{},
	Java:             javaConfig{},
	Zhuge:            zhugeConfig{},
	AMQP:             amqpSetting{},
}

type app struct {
	App              appSetting              `ini:"app"`
	Http             httpSetting             `ini:"http"`
	Database         databaseSetting         `ini:"database"`
	Log              logSetting              `ini:"log"`
	Weapp            wxSetting               `ini:"weapp"`
	MioSubOA         wxSetting               `ini:"mioSubOa"` //绿喵订阅号配置
	MioSrvOA         wxSetting               `ini:"mioSrvOa"` //绿喵服务号配置
	Redis            redisSetting            `ini:"redis"`
	DuiBa            duiBaSetting            `ini:"duiba"`
	OSS              ossSetting              `ini:"oss"`
	AMQP             amqpSetting             `ini:"amqp"`
	BaiDu            baiDuSetting            `ini:"baidu"`
	BaiDuImageSearch baiDuImageSearchSetting `ini:"baiduImageSearch"`
	Java             javaConfig              `ini:"java"`
	Zhuge            zhugeConfig             `ini:"zhuge"`
}
type appSetting struct {
	TokenKey string
	Domain   string
	Debug    bool
	//prod dev local
	Env string
}
type httpSetting struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
	Throttle     string
}
type databaseSetting struct {
	Type         string
	Host         string
	Port         int
	UserName     string
	Password     string
	Database     string
	TablePrefix  string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  int
	LogLevel     string
}
type logSetting struct {
	Level   string
	MaxSize int
}
type wxSetting struct {
	AppId  string
	Secret string
}
type redisSetting struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type duiBaSetting struct {
	AppKey    string
	AppSecret string
}
type ossSetting struct {
	CdnDomain    string
	Endpoint     string
	AccessKey    string
	AccessSecret string
	BasePath     string
}
type amqpSetting struct {
	Url string
}
type baiDuSetting struct {
	AppKey    string
	AppSecret string
}
type baiDuImageSearchSetting struct {
	AppKey    string
	AppSecret string
}
type javaConfig struct {
	JavaLoginUrl string `binding:"required"`
	JavaWhoAmi   string `binding:"required"`
}
type zhugeConfig struct {
	AppKey    string
	AppSecret string
}

func FindOaSetting(source entity.UserSource) wxSetting {
	switch source {
	case entity.UserSourceMioSrvOA:
		return Config.MioSrvOA
	case entity.UserSourceMioSubOA:
		return Config.MioSubOA
	}
	return wxSetting{}
}
