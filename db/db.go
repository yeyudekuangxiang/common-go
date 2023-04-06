package db

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"strings"
	"time"
)

type Config struct {
	Type         string // mysql postgres
	Host         string
	UserName     string
	Password     string
	Database     string
	Port         int
	TablePrefix  string //前缀
	MaxOpenConns int    //最大连接数 <=0表示不限制连接数
	MaxIdleConns int    //最大空闲数 <=0表示不保留空闲连接
	MaxIdleTime  int    //连接最大空闲时间 合理的设置可以避免高并发后长时间占用过多的链接  <=0表示永远可用(单位秒)
	MaxLifetime  int    //连接可重用时间(必须小于数据库wait_timeout) <=0表示永远可用(单位秒)
	Logger       logger.Interface
}
type ReplicasConfig struct {
	List              []ReplicaDBConfig
	MaxOpenConns      int //最大连接数 <=0表示不限制连接数
	MaxIdleConns      int //最大空闲数 <=0表示不保留空闲连接
	MaxIdleTime       int //连接最大空闲时间 合理的设置可以避免高并发后长时间占用过多的链接  <=0表示永远可用(单位秒)
	MaxLifetime       int //连接可重用时间(必须小于数据库wait_timeout) <=0表示永远可用(单位秒)
	TraceResolverMode bool
}
type ReplicaDBConfig struct {
	Host     string
	UserName string
	Password string
	Database string
	Port     int
}
type PoolConfig struct {
	MaxOpenConns int `json:",optional"` //最大连接数 <=0表示不限制连接数
	MaxIdleConns int `json:",optional"` //最大空闲数 <=0表示不保留空闲连接
	MaxIdleTime  int `json:",optional"` //连接最大空闲时间 合理的设置可以避免高并发后长时间占用过多的链接  <=0表示永远可用(单位秒)
	MaxLifetime  int `json:",optional"` //连接可重用时间(必须小于数据库wait_timeout) <=0表示永远可用(单位秒)
}

func NewDB(conf Config, opts ...Option) (*gorm.DB, error) {
	switch strings.ToLower(conf.Type) {
	case "mysql":
		return NewMysqlDB(conf, opts...)
	case "postgres":
		return NewPostgresDB(conf, opts...)
	default:
		return NewMysqlDB(conf, opts...)
	}
}

// NewMysqlDB 创建Mysql数据库链接
func NewMysqlDB(conf Config, opts ...Option) (*gorm.DB, error) {
	options := initOptions(opts...)
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.UserName,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   conf.TablePrefix, // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true,             // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, "open db failed")
	}
	sqlDb, err := db.DB()
	if err == nil {
		sqlDb.SetMaxIdleConns(conf.MaxIdleConns)
		sqlDb.SetMaxOpenConns(conf.MaxOpenConns)
		sqlDb.SetConnMaxIdleTime(time.Duration(conf.MaxIdleTime) * time.Second)
		sqlDb.SetConnMaxLifetime(time.Duration(conf.MaxLifetime) * time.Second)
	}
	//配置日志打印
	if conf.Logger != nil {
		db.Logger = conf.Logger
	}
	db.Logger = conf.Logger
	if options.replicas != nil {
		err = db.Use(newMysqlDBResolver(*options.replicas))
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}
func newMysqlDBResolver(c ReplicasConfig) *dbresolver.DBResolver {
	if len(c.List) == 0 {
		panic("replica db count must gt 0")
	}
	replicas := make([]gorm.Dialector, 0)
	for _, r := range c.List {
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			r.UserName,
			r.Password,
			r.Host,
			r.Port,
			r.Database,
		)
		replicas = append(replicas, mysql.Open(dsn))
	}
	return dbresolver.Register(dbresolver.Config{
		Replicas:          replicas,
		TraceResolverMode: c.TraceResolverMode,
	}).
		SetMaxIdleConns(c.MaxIdleConns).
		SetMaxOpenConns(c.MaxOpenConns).
		SetConnMaxIdleTime(time.Duration(c.MaxIdleTime) * time.Second).
		SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)
}
func newPostgresDBResolver(c ReplicasConfig) *dbresolver.DBResolver {
	if len(c.List) == 0 {
		panic("replica db count must gt 0")
	}
	replicas := make([]gorm.Dialector, 0)
	for _, r := range c.List {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
			r.Host,
			r.UserName,
			r.Password,
			r.Database,
			r.Port,
		)
		replicas = append(replicas, postgres.Open(dsn))
	}
	return dbresolver.Register(dbresolver.Config{
		Replicas:          replicas,
		TraceResolverMode: c.TraceResolverMode,
	}).SetMaxIdleConns(c.MaxIdleConns).
		SetMaxOpenConns(c.MaxOpenConns).
		SetConnMaxIdleTime(time.Duration(c.MaxIdleTime) * time.Second).
		SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)
}

type Options struct {
	replicas *ReplicasConfig
}
type Option func(options *Options)

func WithReplicas(replicas ReplicasConfig) Option {
	return func(options *Options) {
		options.replicas = &replicas
	}
}
func initOptions(opts ...Option) Options {
	option := Options{}
	for _, opt := range opts {
		opt(&option)
	}
	return option
}

//NewPostgresDB 创建PostgreSQL数据库链接
func NewPostgresDB(conf Config, opts ...Option) (*gorm.DB, error) {
	options := initOptions(opts...)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		conf.Host,
		conf.UserName,
		conf.Password,
		conf.Database,
		conf.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   conf.TablePrefix, // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true,             // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, "open db failed")
	}
	sqlDb, err := db.DB()
	if err == nil {
		sqlDb.SetMaxIdleConns(conf.MaxIdleConns)
		sqlDb.SetMaxOpenConns(conf.MaxOpenConns)
		sqlDb.SetConnMaxIdleTime(time.Duration(conf.MaxIdleTime) * time.Second)
		sqlDb.SetConnMaxLifetime(time.Duration(conf.MaxLifetime) * time.Second)
	}
	//配置日志打印
	if conf.Logger != nil {
		db.Logger = conf.Logger
	}
	if options.replicas != nil {
		err = db.Use(newPostgresDBResolver(*options.replicas))
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}
