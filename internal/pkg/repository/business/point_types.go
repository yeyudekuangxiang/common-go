package business

import (
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/model/entity/business"
	"time"
)

type GetPointTransactionListBy struct {
	UserId    int64
	StartTime time.Time
	EndTime   time.Time
	OrderBy   entity.OrderByList
	Type      business.PointType
}
