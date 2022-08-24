package event

import (
	"database/sql"
	"mio/internal/pkg/model/entity"
)

type FindEventBy struct {
	ProductItemId string
	EventId       string
	Active        sql.NullBool
}

type GetEventPageListBy struct {
	Limit           int
	IsShow          int //1显示 2不显示
	Offset          int
	EventCategoryId string
	OrderBy         entity.OrderByList
}
type GetEventListBy struct {
	IsShow          int //1显示 2不显示
	EventCategoryId string
	OrderBy         entity.OrderByList
	Active          sql.NullBool
}
type GetEventCategoryListBy struct {
	OrderBy entity.OrderByList
	Active  sql.NullBool
}
type GetEventDetailListBy struct {
	EventId string
}
type GetEventRuleListBy struct {
	EventId string
}
type GetParticipationPageListBy struct {
	EventId string
	Limit   int
	Offset  int
	OrderBy entity.OrderByList
}
