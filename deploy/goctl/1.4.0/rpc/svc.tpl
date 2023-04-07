package svc

import {{.imports}}
import (
	"{{.projectPath}}/common/tool/db"
	"gitlab.miotech.com/miotech-application/backend/common-go/logger/aliyunzero"
	"gitlab.miotech.com/miotech-application/backend/common-go/logger/gormzero"
	"context"
	"log"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	gormDB, err := db.NewDB(c.Database,gormzero.NewLogger(aliyunzero.WithFiledLogger(logx.WithContext(context.Background()), aliyunzero.LogTopicDatabase), time.Second))
	if err != nil {
		log.Panicf("数据库连接失败 %+v %v", c.Database, err)
	}
	return &ServiceContext{
		Config:c,
	}
}

