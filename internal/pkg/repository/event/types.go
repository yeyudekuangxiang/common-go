package event

import "mio/internal/pkg/model/entity"

type FindEventBy struct {
	ProductItemId string
	EventId       string
}

type GetEventPageListBy struct {
	Limit           int
	Offset          int
	EventCategoryId string
	OrderBy         entity.OrderByList
}
type GetEventListBy struct {
	EventCategoryId string
	OrderBy         entity.OrderByList
}
type GetEventCategoryListBy struct {
	OrderBy entity.OrderByList
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
