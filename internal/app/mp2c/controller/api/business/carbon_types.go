package business

import (
	"mio/internal/pkg/model/entity/business"
	"time"
)

type GetCarbonCreditLogInfoListForm struct {
	Date time.Time `json:"date" form:"date" binding:"required" alias:"月份" time_format:"2006-01" time_utc:"false" time_location:"Asia/Shanghai"`
}
type ChangeUserLikeStatusForm struct {
	Pid        int64                   `json:"pid"`
	ObjectType business.RankObjectType `json:"objectType" form:"objectType" binding:"oneof=user department" alias:"对象类型"`
	DateType   business.RankDateType   `json:"dateType" form:"dateType" binding:"oneof=day week month" alias:"排行榜类型"`
}

type CarbonCollectEvCarForm struct {
	Electricity int64 `json:"electricity" form:"electricity" required:"gte=0" alias:"电量"`
}
type CarbonCollectOnlineMeetingForm struct {
	Duration float64 `json:"duration" form:"duration" required:"gt=0" alias:"会议时长"`
}
type CarbonCollectSaveWaterElectricityForm struct {
	Water       int64 `json:"water" form:"water" `
	Electricity int64 `json:"electricity" form:"electricity"`
}
type CarbonCollectPublicTransportForm struct {
	Bus   int64
	Metro int64
}
