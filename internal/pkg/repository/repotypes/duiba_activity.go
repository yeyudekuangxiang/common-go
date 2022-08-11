package repotypes

import "mio/internal/pkg/model/entity"

type GetDuiBaActivityExistDO struct {
	ActivityId string
	Id         int64
	NotId      int64
}
type GetDuiBaActivityPageDO struct {
	Cid        int64
	Type       entity.DuiBaActivityType
	ActivityId string
	Name       string
	Offset     int                `json:"offset"`
	Limit      int                `json:"limit"` //limit为0时不限制数量
	OrderBy    entity.OrderByList `json:"orderBy"`
	Statue     entity.DuiBaActivityStatus
}
