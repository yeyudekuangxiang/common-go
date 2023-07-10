package {{.pkg}}

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/syncx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

var (
    singleFlights = syncx.NewSingleFlight()
    stats = cache.NewStat("sqlc")
)


// see doc/sql-cache.md
const cacheSafeGapBetweenIndexAndPrimary = time.Second * 5

type Options struct {
	//查询数据是否跳过缓存直接从数据库中查询
    skipFindCache bool
}

type option func(db *gorm.DB, options Options) (*gorm.DB, Options)

//ForUpdate 查询一条数据用于更新(加行级锁) 在查询方法中使用会跳过缓存直接从数据库查询
func ForUpdate() option {
	return func(db *gorm.DB, options Options) (*gorm.DB, Options) {
		options.skipFindCache = true
		return db.Clauses(clause.Locking{Strength: "UPDATE"}), options
	}
}

// FindSkipCache 查询跳过缓存直接从数据库查询
func FindSkipCache() option {
	return func(db *gorm.DB, options Options) (*gorm.DB, Options) {
		options.skipFindCache = true
		return db, options
	}
}

func initOptions(db *gorm.DB, options Options, opts []option) (*gorm.DB, Options) {
	for _, opt := range opts {
		db, options = opt(db, options)
	}
	return db, options
}
