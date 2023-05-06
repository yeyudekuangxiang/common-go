package initialize

import (
	"github.com/zeromicro/go-zero/core/conf"
	"gopkg.in/ini.v1"
	"log"
	"mio/config"
	"mio/internal/pkg/service/duiba"
	"mio/internal/pkg/service/oss"
)

func InitIni(source interface{}) {
	log.Println("初始化配置文件...")
	f, err := ini.Load(source)
	if err != nil {
		log.Panic(err)
	}
	err = f.MapTo(&config.Config)
	if err != nil {
		log.Panic(err)
	}
	log.Println("初始化配置文件成功")
	log.Printf("%+v\n", config.Config)
	afterInitIni()
}
func InitYaml(source string) {
	err := conf.Load(source, &config.Config)
	if err != nil {
		log.Panic(err)
	}
	log.Println("初始化配置文件成功")
	log.Printf("%+v\n", config.Config)
	afterInitIni()
}
func afterInitIni() {
	duiba.InitDefaultDuibaService()
	oss.InitDefaultOssService()
}
