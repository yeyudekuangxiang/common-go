package srv_types

import (
	"mio/internal/pkg/model/entity/event"
)

type FindBadgeParam struct {
	OrderId string
}
type UploadBadgeResult struct {
	EventTemplateType event.EventTemplateType `json:"eventTemplateType"`
	TemplateSetting   map[string]interface{}  `json:"templateSetting"`
	BadgeInfo         UploadBadgeInfo         `json:"badgeInfo"`
}
type UploadBadgeInfo struct {
	UploadCode    string `json:"uploadCode"`
	Username      string `json:"username"`
	Time          string `json:"time"`
	CertificateNo string `json:"certificateNo"`
}
