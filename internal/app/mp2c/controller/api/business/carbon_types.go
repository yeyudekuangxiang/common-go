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
	Electricity float64 `json:"electricity" form:"electricity" binding:"gt=0,lte=30" alias:"电量"`
}
type CarbonCollectOnlineMeetingForm struct {
	OneCityDuration  float64 `json:"oneCityDuration" form:"oneCityDuration" binding:"required_if=ManyCityDuration 0,lte=3" alias:"同城会议"`
	ManyCityDuration float64 `json:"manyCityDuration" form:"manyCityDuration" binding:"required_if=OneCityDuration 0,lte=3" alias:"异地会议"`
}
type CarbonCollectSaveWaterElectricityForm struct {
	Water       int64 `json:"water" form:"water" binding:"required_if=Electricity 0" alias:"水量"`
	Electricity int64 `json:"electricity" form:"electricity" binding:"required_if=Water 0" alias:"电量"`
}
type CarbonCollectPublicTransportForm struct {
	Bus   float64 `json:"bus" form:"bus" binding:"required_if=Metro 0 Walk 0 Bike 0,lte=30" alias:"公交"`
	Metro float64 `json:"metro" form:"metro" binding:"required_if=Bus 0 Walk 0 Bike 0,lte=30" alias:"地铁"`
	Walk  float64 `json:"walk" form:"walk" binding:"required_if=Bus 0 Metro 0 Bike 0,lte=10" alias:"步行"`
	Bike  float64 `json:"bike" form:"bike" binding:"required_if=Bus 0 Metro 0 Walk 0,lte=10" alias:"骑行"`
}
type GetCarbonCreditLogSortedListForm struct {
	StartTime time.Time `json:"startTime" form:"startTime" binding:"required" alias:"开始时间" time_format:"2006-01-02" time_utc:"false" time_location:"Asia/Shanghai"`
	EndTime   time.Time `json:"endTime" form:"endTime" binding:"required" alias:"结束时间" time_format:"2006-01-02" time_utc:"false" time_location:"Asia/Shanghai"`
}

type CarbonCollectOEPForm struct {
	Photo string `json:"photo" form:"photo" binding:"required" alias:"凭证图片"`
}
type CarbonGreenBusinessTripForm struct {
	Type  string `json:"type" form:"type" binding:"oneof=hsr train plane" alias:"交通工具"`
	From  string `json:"from" form:"from" binding:"required" alias:"出发地"`
	To    string `json:"to" form:"to" binding:"required" alias:"目的地"`
	Photo string `json:"photo" form:"photo" binding:"required" alias:"凭证图片"`
}