package event

import (
	"database/sql"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	eevent "mio/internal/pkg/model/entity/event"
)

type ParticipateEventParam struct {
	ProductItemId string
	Count         int
}
type FindEventParam struct {
	ProductItemId string
	EventId       string
	Active        sql.NullBool
}
type GetEventDetailListParam struct {
	EventId string
}
type GetEventRuleListParam struct {
	EventId string
}
type EventShortInfo struct {
	EventId           string                   `json:"eventId" `
	EventTemplateType eevent.EventTemplateType `json:"eventTemplateType"`
	Title             string                   `json:"title" `
	Subtitle          string                   `json:"subtitle" `
	CoverImageUrl     string                   `json:"coverImageUrl" `
	Cost              int                      `json:"cost" `
}

type EventFullInfo struct {
	EventId               string                                     `json:"eventId"`
	EventTemplateType     eevent.EventTemplateType                   `json:"eventTemplateType"`
	Title                 string                                     `json:"title"`
	SubTitle              string                                     `json:"subTitle"`
	CoverImageUrl         string                                     `json:"coverImageUrl"`
	StartTime             model.Time                                 `json:"startTime"`
	EndTime               model.Time                                 `json:"endTime"`
	ParticipationCount    int                                        `json:"participationCount"`
	ParticipationSubtitle string                                     `json:"participationSubtitle"`
	Tags                  []string                                   `json:"tag"`
	TemplateSetting       map[string]eevent.EventTemplateSettingInfo `json:"templateSetting"`
	ParticipationList     []ParticipationInfo                        `json:"participationList"`
	EventDetail           string                                     `json:"eventDetail"`
	EventRule             string                                     `json:"eventRule"`
	Cost                  int                                        `json:"cost"`
}
type ParticipationInfo struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Message  string `json:"message"`
	TimeStr  string `json:"timeStr"`
}
type GetParticipationPageListParam struct {
	EventId string
	Limit   int
	Offset  int
	OrderBy entity.OrderByList
}
type UserParticipationCount struct {
	UserId int
	Count  int
}
type GetEventListParam struct {
	EventCategoryId string
	OrderBy         entity.OrderByList
	Active          sql.NullBool
}

type GetEventCategoryListParam struct {
	OrderBy entity.OrderByList
	Active  sql.NullBool
}
