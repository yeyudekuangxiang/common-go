package repotypes

import "mio/internal/pkg/model/entity"

type GetPageFullOrderDO struct {
	Openid      string
	OrderSource entity.OrderSource
	Limit       int
	Offset      int
}

type GetOrderTotalByItemIdDO struct {
	Openid      string
	ItemIdSlice []string
	StartTime   string
	EndTime     string
}
