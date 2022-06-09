package business

import (
	"mio/internal/pkg/model/entity/business"
	"time"
)

type GetCarbonCreditLogInfoListForm struct {
	Date time.Time `json:"date" form:"date" binding:"required" alias:"月份" time_format:"2006-01" time_utc:"false" time_location:"Asia/Shanghai"`
}
type ChangeUserRankLikeStatusForm struct {
	Uid      string                `json:"uid" form:"uid" binding:"required" alias:"用户编号"`
	DateType business.RankDateType `json:"dateType" form:"dateType" binding:"oneof=day week month" alias:"排行榜类型"`
}
type ChangeDepartmentRankLikeStatusForm struct {
	DepartmentId int64                 `json:"departmentId" form:"departmentId" binding:"required" alias:"部门ID"`
	DateType     business.RankDateType `json:"dateType" form:"dateType" binding:"oneof=day week month" alias:"排行榜类型"`
}
type CarbonCollectEvCarForm struct {
	Electricity float64 `json:"electricity" form:"electricity" required:"gte=0" alias:"电量"`
}
type CarbonCollectOnlineMeetingForm struct {
	OnlineDuration  float64 `json:"duration" form:"duration" required:"gt=0" alias:"异地会议"`
	OfflineDuration float64 `json:"offlineDuration" form:"offlineDuration" required:"gt=0" alias:"同城会议"`
}
type CarbonCollectSaveWaterElectricityForm struct {
	Water       int64 `json:"water" form:"water" `
	Electricity int64 `json:"electricity" form:"electricity"`
}
type CarbonCollectPublicTransportForm struct {
	Bus   int64 `json:"bus" form:"bus" binding:"gte=0" alias:"公交"`
	Metro int64 `json:"metro" form:"metro" binding:"gte=0" alias:"地铁"`
}
