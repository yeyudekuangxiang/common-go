package business

import (
	"github.com/shopspring/decimal"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/model/entity/business"
	"time"
)

type GetPointTransactionListParam struct {
	UserId    int64
	StartTime time.Time
	EndTime   time.Time
	OrderBy   entity.OrderByList
	Type      business.PointType
}

type PointTransactionInfo struct {
	ID       int64              `json:"id"`
	Type     business.PointType `json:"type"`
	TypeText string             `json:"typeText"`
	TimeStr  string             `json:"timeStr"`
	Value    int                `json:"value"`
}
type GetPointTransactionInfoListParam struct {
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
