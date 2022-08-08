package repotypes

import (
	"mio/internal/pkg/model/entity"
	"time"
)

type GetCarbonTransactionCountLimitDO struct {
	OpenId    string
	StartTime time.Time
	EndTime   time.Time
	OrderBy   entity.OrderByList
	Type      entity.CarbonTransactionType
	VDate     time.Time
}

type FindCarbonTransactionCountLimitFindByDO struct {
	OpenId string
	Type   entity.CarbonTransactionType
	VDate  time.Time
}
