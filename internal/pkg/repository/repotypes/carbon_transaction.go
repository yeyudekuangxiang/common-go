package repotypes

import (
	"mio/internal/pkg/model/entity"
)

type GetCarbonTransactionListByDO struct {
	Uid       int64
	OpenId    string
	StartTime string
	EndTime   string
	OrderBy   entity.OrderByList
	Type      entity.CarbonTransactionType
}

type GetCarbonTransactionListBy struct {
	Type   entity.CarbonTransactionType
	Sum    float64
	UserId int64
	Openid string
}
