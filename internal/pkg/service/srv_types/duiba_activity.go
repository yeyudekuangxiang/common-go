package srv_types

import "mio/internal/pkg/model/entity"

type CreateDuiBaActivityDTO struct {
	Name        string
	Cid         int64
	Type        entity.DuiBaActivityType
	IsShare     entity.DuiBaActivityIsShare
	ActivityUrl string
	ActivityId  string
}

type UpdateDuiBaActivityDTO struct {
	Id          int64
	Name        string
	Cid         int64
	Type        entity.DuiBaActivityType
	IsShare     entity.DuiBaActivityIsShare
	ActivityUrl string
	ActivityId  string
}

type DeleteDuiBaActivityDTO struct {
	Id int64
}

type GetPageDuiBaActivityDTO struct {
	Cid        int64
	Type       entity.DuiBaActivityType
	Status     entity.DuiBaActivityStatus
	ActivityId string
	Name       string
	Offset     int                `json:"offset"`
	Limit      int                `json:"limit"` //limit为0时不限制数量
	OrderBy    entity.OrderByList `json:"orderBy"`
}
