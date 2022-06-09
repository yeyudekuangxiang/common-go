package business

import (
	"github.com/shopspring/decimal"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/model/entity/business"
	"time"
)

type GetPointLogListParam struct {
	UserId    int64
	StartTime time.Time
	EndTime   time.Time
	OrderBy   entity.OrderByList
	Type      business.PointType
}

type PointLogInfo struct {
	ID       int64              `json:"id"`
	Type     business.PointType `json:"type"`
	TypeText string             `json:"typeText"`
	TimeStr  string             `json:"timeStr"`
	Value    int                `json:"value"`
}
type GetPointLogInfoListParam struct {
	UserId    int64
	StartTime time.Time
	EndTime   time.Time
}

type GetCarbonCreditLogListParam struct {
	UserId    int64
	StartTime time.Time
	EndTime   time.Time
	OrderBy   entity.OrderByList
	Type      business.CarbonType
}
type CarbonCreditLogInfo struct {
	ID       int64               `json:"id"`
	Type     business.CarbonType `json:"type"`
	TypeText string              `json:"typeText"`
	TimeStr  string              `json:"timeStr"`
	Value    decimal.Decimal     `json:"value"`
}
type GetCarbonCreditLogInfoListParam struct {
	UserId    int64
	StartTime time.Time
	EndTime   time.Time
}
type CreateCarbonCreditLogParam struct {
	TransactionId string
	UserId        int64
	Type          business.CarbonType
	Value         decimal.Decimal
	Info          business.CarbonTypeInfo
}
type FindPointLimitLogParam struct {
	TimePoint time.Time
	Type      business.PointType
	UserId    int64
}
type UpdateOrCreatePointLimitLogParam struct {
	Type            business.PointType
	UserId          int64
	AddCurrentValue int //增加的积分数量
	TimePoint       time.Time
}
type CreatePointLogParam struct {
	TransactionId string
	UserId        int64
	Type          business.PointType
	Value         int
	Info          business.PointTypeInfo
	OrderId       string
}

type SendPointParam struct {
	UserId        int64
	AddPoint      int
	Info          business.PointTypeInfo
	Type          business.PointType
	OrderId       string
	TransactionId string
}
type createOrUpdatePointParam struct {
	UserId   int64
	AddPoint int
}
