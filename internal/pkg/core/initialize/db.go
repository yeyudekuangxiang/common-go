package initialize

import (
	mzap "gitlab.miotech.com/miotech-application/backend/common-go/logger/zap"
	"gorm.io/gorm/logger"
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/pkg/db"
)

//Silent Error Warn Info
var gormLevelMap = map[string]logger.LogLevel{
	"silent": logger.Silent,
	"error":  logger.Error,
	"warn":   logger.Warn,
	"info":   logger.Info,
	"":       logger.Info,
}

func InitDB() {
	log.Println("初始化数据库连接...")
	dbc := config.Config.Database

	zlogger := app.OriginLogger.With(mzap.LogDatabase)
	conf := db.Config{
		Type:         dbc.Type,
		Host:         dbc.Host,
		UserName:     dbc.UserName,
		Password:     dbc.Password,
		Database:     dbc.Database,
		Port:         dbc.Port,
		TablePrefix:  dbc.TablePrefix,
		MaxOpenConns: dbc.MaxOpenConns,
		MaxIdleConns: dbc.MaxIdleConns,
		MaxLifetime:  dbc.MaxLifetime,
		Logger:       mzap.NewGormLogger(zlogger.Sugar()).LogMode(gormLevelMap[dbc.LogLevel]),
	}

	gormDb, err := db.NewDB(conf)
	if err != nil {
		log.Panic(err)
	}
	*app.DB = *gormDb
	log.Println("初始化数据库连接成功")
}

func InitBusinessDB() {
	log.Println("初始化企业版数据库连接...")
	dbc := config.Config.DatabaseBusiness

	zlogger := app.OriginLogger.With(mzap.LogDatabase)
	conf := db.Config{
		Type:         dbc.Type,
		Host:         dbc.Host,
		UserName:     dbc.UserName,
		Password:     dbc.Password,
		Database:     dbc.Database,
		Port:         dbc.Port,
		TablePrefix:  dbc.TablePrefix,
		MaxOpenConns: dbc.MaxOpenConns,
		MaxIdleConns: dbc.MaxIdleConns,
		MaxLifetime:  dbc.MaxLifetime,
		Logger:       mzap.NewGormLogger(zlogger.Sugar()).LogMode(gormLevelMap[dbc.LogLevel]),
	}

	gormDb, err := db.NewDB(conf)
	if err != nil {
		log.Panic(err)
	}
	*app.BusinessDB = *gormDb
	log.Println("初始化企业版数据库连接成功")
}
