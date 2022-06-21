package event

import (
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
}
type GetEventDetailListParam struct {
	EventId string
}
type GetEventRuleListParam struct {
	EventId string
}
type EventFullInfo struct {
	EventId               string                          `json:"eventId"`
	EventTemplateType     eevent.EventTemplateType        `json:"eventTemplateType"`
	Title                 string                          `json:"title"`
	SubTitle              string                          `json:"subTitle"`
	CoverImageUrl         string                          `json:"coverImageUrl"`
	StartTime             model.Time                      `json:"startTime"`
	EndTime               model.Time                      `json:"endTime"`
	ParticipationCount    int                             `json:"participationCount"`
	ParticipationTitle    string                          `json:"participationTitle"`
	ParticipationSubtitle string                          `json:"participationSubtitle"`
	Tags                  []string                        `json:"tag"`
	TemplateSetting       eevent.EventTemplateSettingInfo `json:"templateSetting"`
	ParticipationList     []ParticipationInfo             `json:"participationList"`
	EventDetail           string                          `json:"eventDetail"`
	EventRule             string                          `json:"eventRule"`
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
}

type GetEventCategoryListParam struct {
	OrderBy entity.OrderByList
}
