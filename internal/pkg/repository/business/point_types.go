package business

import (
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/model/entity/business"
	"time"
)

type GetPointLogListBy struct {
	UserId    int64
	StartTime time.Time
	EndTime   time.Time
	OrderBy   entity.OrderByList
	Type      business.PointType
}
type FindPointLimitLogBy struct {
	TimePoint time.Time
	Type      business.PointType
	UserId    int64
}
type FindCarbonCreditsLimitLogBy struct {
	TimePoint time.Time
	Type      business.CarbonType
	UserId    int64
}
type FindPointBy struct {
	UserId int64
}
