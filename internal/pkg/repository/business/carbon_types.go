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

type GetActualDepartmentCarbonRankBy struct {
	StartTime time.Time
	EndTime   time.Time
	CompanyId int
	Limit     int
	Offset    int
}

type GetActualUserCarbonRankBy struct {
	StartTime time.Time
	EndTime   time.Time
	CompanyId int
	Limit     int
	Offset    int
}
type GetCarbonRankBy struct {
	TimePoint  time.Time
	DateType   ebusiness.RankDateType
	ObjectType ebusiness.RankObjectType
	CompanyId  int
	Limit      int
	Offset     int
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
type FindCarbonRankBy struct {
	Pid        int64
	ObjectType ebusiness.RankObjectType
	DateType   ebusiness.RankDateType
	TimePoint  time.Time
}
type FindCompanyCarbonSceneBy struct {
	CompanyId     int
	CarbonSceneId int
}
