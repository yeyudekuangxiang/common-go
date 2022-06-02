package business

import (
	"mio/internal/pkg/model/entity"
	ebusiness "mio/internal/pkg/model/entity/business"
	"time"
)

type GetCarbonCreditsLogListBy struct {
	UserId    int64
	StartTime time.Time
	EndTime   time.Time
	OrderBy   entity.OrderByList
	Type      ebusiness.CarbonType
}

type GetUserCarbonRankBy struct {
	StartTime time.Time
	EndTime   time.Time
	UserId    int64
	CompanyId int
	Limit     int
	Offset    int
}
type GetDepartmentCarbonRankBy struct {
	StartTime    time.Time
	EndTime      time.Time
	DepartmentId int
	CompanyId    int
	Limit        int
	Offset       int
}

type GetCarbonRankLikeNumListBy struct {
	PIds       []int64
	ObjectType ebusiness.RankObjectType
	DateType   ebusiness.RankDateType
	TimePoint  time.Time
}
type FindCarbonRankLikeNumBy struct {
	Pid        int64
	ObjectType ebusiness.RankObjectType
	DateType   ebusiness.RankDateType
	TimePoint  time.Time
}
type FindCarbonRankLikeLogBy struct {
	Pid        int64
	UserId     int64
	ObjectType ebusiness.RankObjectType
	DateType   ebusiness.RankDateType
	TimePoint  time.Time
}
type GetCarbonRankLikeLogListBy struct {
	PIds       []int64
	UserId     int64
	ObjectType ebusiness.RankObjectType
	DateType   ebusiness.RankDateType
	TimePoint  time.Time
}
