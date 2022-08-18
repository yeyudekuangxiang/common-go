package businesstypes

import (
	"mio/internal/pkg/model/entity/business"
	"time"
)

type GetPointRecordListForm struct {
	Date time.Time `json:"date" form:"date" binding:"required" alias:"月份" time_format:"2006-01" time_utc:"false" time_location:"Asia/Shanghai"`
}
type PointLogInfo struct {
	ID       int64              `json:"id"`
	Type     business.PointType `json:"type"`
	TypeText string             `json:"typeText"`
	TimeStr  string             `json:"timeStr"`
	Value    int                `json:"value"`
}
