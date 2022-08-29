package repotypes

import (
	"mio/internal/pkg/model/entity"
)

type GetCarbonTransactionDayGetListDO struct {
	UserId    int64
	StartTime string
	EndTime   string
	OrderBy   entity.OrderByList
}
