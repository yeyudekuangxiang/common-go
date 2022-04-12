package app

import (
	"github.com/go-redis/redis/v8"
	"github.com/medivhzhan/weapp/v3"
	"gopkg.in/ini.v1"
	"gorm.io/gorm"
	"mio/pkg/zap"
)

var (
	DB     = new(gorm.DB)
	Ini    *ini.File
	Logger = zap.DefaultLogger().Sugar()
	Weapp  = new(weapp.Client)
	Redis  = new(redis.Client)
)
