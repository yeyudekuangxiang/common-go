package db

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"log"
	"os"
	"testing"
	"time"
)

func TestNewDBWithReplicas(t *testing.T) {
	db, err := NewDB(Config{
		Type:         "postgres",
		Host:         "pgm-uf680kf1780dfqvt.pg.rds.aliyuncs.com",
		UserName:     "miniprogram",
		Password:     "f8N27Tj3Is3ZIbU",
		Database:     "miniprogram",
		Port:         5432,
		TablePrefix:  "",
		MaxOpenConns: 10,
		MaxIdleConns: 2,
		MaxIdleTime:  50,
		MaxLifetime:  100,
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		}),
	}, WithReplicas(ReplicasConfig{
		List: []ReplicaDBConfig{
			{
				Host:     "pgm-uf680kf1780dfqvt.pg.rds.aliyuncs.com",
				UserName: "miniprogram",
				Password: "f8N27Tj3Is3ZIbU",
				Database: "miniprogram",
				Port:     5432,
			},
		},
		MaxOpenConns:      10,
		MaxIdleConns:      2,
		MaxIdleTime:       50,
		MaxLifetime:       100,
		TraceResolverMode: true,
	}))
	assert.Equal(t, nil, err)

	type User struct {
		ID int64
	}
	u := User{}
	err = db.Clauses(dbresolver.Write).Raw("SELECT id FROM \"user\" WHERE id = ?", 525946).Scan(&u).Error
	assert.Equal(t, nil, err)
}

func TestNewDB(t *testing.T) {
	db, err := NewDB(Config{
		Type:         "postgres",
		Host:         "pgm-uf680kf1780dfqvt.pg.rds.aliyuncs.com",
		UserName:     "miniprogram",
		Password:     "f8N27Tj3Is3ZIbU",
		Database:     "miniprogram",
		Port:         5432,
		TablePrefix:  "",
		MaxOpenConns: 10,
		MaxIdleConns: 2,
		MaxIdleTime:  50,
		MaxLifetime:  100,
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		}),
	})
	assert.Equal(t, nil, err)

	type User struct {
		ID int64
	}
	u := User{}
	err = db.Raw("SELECT id FROM \"user\" WHERE id = ?", 525946).Scan(&u).Error
	assert.Equal(t, nil, err)
}
