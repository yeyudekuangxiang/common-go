package event

import (
	eevent "mio/internal/pkg/model/entity/event"
)

type EventCategoryInfo struct {
	EventCategoryId string `json:"eventCategoryId"`
	Title           string `json:"title"`
	ImageUrl        string `json:"imageUrl"`
}
type EventInfo struct {
	EventId           string                   `json:"eventId" gorm:"type:varchar(255);not null;comment:公益活动标识"`
	EventTemplateType eevent.EventTemplateType `json:"eventTemplateType" gorm:"type:varchar(50);not null;default:'';comment:证书模版类型"`
	Title             string                   `json:"title" gorm:"type:varchar(255);not null;comment:公益活动标题"`
	Subtitle          string                   `json:"subtitle" gorm:"type:varchar(255);not null;default:'';comment:公益活动副标题"`
	CoverImageUrl     string                   `json:"coverImageUrl" gorm:"type:varchar(500);not null;comment:项目主图 343 × 200"`
}
