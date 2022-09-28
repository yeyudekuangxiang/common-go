package {{.pkg}}

import (
    "github.com/zeromicro/go-zero/core/stores/sqlx"
    "github.com/zeromicro/go-zero/core/stores/cache"
    "github.com/zeromicro/go-zero/core/syncx"
)

var (
    ErrNotFound = sqlx.ErrNotFound
    singleFlights = syncx.NewSingleFlight()
    stats = cache.NewStat("sqlc")
)

