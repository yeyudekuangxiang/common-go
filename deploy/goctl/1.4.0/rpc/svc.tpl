package svc

import {{.imports}}
import (
	"{{.projectPath}}/common/tool/db"
	"log"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	gormDB, err := db.NewDBWithLogx(c.Database)
	if err != nil {
		log.Panicf("数据库连接失败 %+v %v", c.Database, err)
	}
	return &ServiceContext{
		Config:c,
	}
}

