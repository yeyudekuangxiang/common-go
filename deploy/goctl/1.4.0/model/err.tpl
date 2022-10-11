package {{.pkg}}

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/syncx"
	"gorm.io/gorm"
)

var (
    singleFlights = syncx.NewSingleFlight()
    stats = cache.NewStat("sqlc")
)

type option func(db *gorm.DB) *gorm.DB

//ForUpdate 查询一条数据用于更新(加行级锁)
func ForUpdate() option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Set("gorm:query_option", "for update")
	}
}

func initOptions(db *gorm.DB, opts []option) *gorm.DB {
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}
