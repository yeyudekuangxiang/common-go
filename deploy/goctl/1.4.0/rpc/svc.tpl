package svc

import {{.imports}}

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	gormDB, err := newDB(c.Database)
	if err != nil {
		log.Panicf("数据库连接失败 %+v %v", c.Database, err)
	}
	return &ServiceContext{
		Config:c,
	}
}

func newDB(c config.DbConf)(*gorm.DB,err)
{
    return db.NewDB(db.Config{
		Type:         c.Type,
		Host:         c.Host,
		UserName:     c.UserName,
		Password:     c.Password,
		Database:     c.Database,
		Port:         c.Port,
		TablePrefix:  c.TablePrefix,
		MaxOpenConns: c.MaxOpenConns,
		MaxIdleConns: c.MaxIdleConns,
		MaxLifetime:  c.MaxLifetime,
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: false,       // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,       // 禁用彩色打印
		}),
	})
}