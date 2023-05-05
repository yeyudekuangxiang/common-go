package {{.pkg}}

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/syncx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

var (
    // singleFlights 在缓存中使用,用于查询缓存时只允许一个线程进入查询,其它线程等待
    singleFlights = syncx.NewSingleFlight()
    // stats 统计缓存
    stats = cache.NewStat("sqlc")
)


// cacheSafeGapBetweenIndexAndPrimary 主键缓存比索引缓存多缓存一些时间 避免有索引缓存但没有主键缓存的情况
const cacheSafeGapBetweenIndexAndPrimary = time.Second * 5
// cacheUpdateSync 主从同步间隔时间
const cacheUpdateSync = time.Second * 10

// 更新删除数据后标记名称
const updateTagKey = "dbupdatetag:"

type Options struct {
	// skipFindCache 查询数据是否跳过缓存直接从数据库中查询
    skipFindCache bool
}

type option func(db *gorm.DB, options Options) (*gorm.DB, Options)
type modelOption func(options Options) Options

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

func initMobileOptions(options Options, opts []modelOption) Options {
	for _, opt := range opts {
		options = opt(options)
	}
	return options
}

func ModelOptionSkipCache() modelOption {
	return func(options Options) Options {
		options.skipFindCache = true
		return options
	}
}