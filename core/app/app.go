package app

import (
	"github.com/medivhzhan/weapp/v3"
	"gopkg.in/ini.v1"
	"gorm.io/gorm"
	"mio/internal/zap"
	"net/http"
)

var (
	DB     = new(gorm.DB)
	Ini    *ini.File
	Logger = zap.DefaultLogger().Sugar()
	Server *http.Server
	Weapp  *weapp.Client
)
