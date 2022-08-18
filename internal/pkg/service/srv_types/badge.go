package srv_types

import (
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity/event"
)

type FindBadgeParam struct {
	OrderId string
}
type UploadOldBadgeResult struct {
	CreateTime        model.Time `json:"createTime"`
	EventTemplateType event.EventTemplateType
	TemplateSetting   map[string]interface{} `json:"templateSetting"`
	UploadCode        string                 `json:"uploadCode"`
}
