package business

import (
	"github.com/shopspring/decimal"
	"mio/internal/pkg/model/entity/business"
	"time"
)

type GetCarbonRankLikeLogListParam struct {
	PIds       []int64
	UserId     int64
	DateType   business.RankDateType
	ObjectType business.RankObjectType
	TimePoint  time.Time
}
type CarbonRankLikeLogParam struct {
	Pid        int64
	UserId     int64
	DateType   business.RankDateType
	ObjectType business.RankObjectType
	TimePoint  time.Time
}
type GetCarbonRankLikeNumListParam struct {
	PIds       []int64
	ObjectType business.RankObjectType
	DateType   business.RankDateType
	TimePoint  time.Time
}
type FindCarbonRankLikeNumParam struct {
	Pid        int64
	ObjectType business.RankObjectType
	DateType   business.RankDateType
	TimePoint  time.Time
}

type CreateCarbonRankLikeNumParam struct {
	Pid        int64
	ObjectType business.RankObjectType
	DateType   business.RankDateType
	TimePoint  time.Time
}
type UpdateCarbonRankLikeNumParam struct {
	Pid        int64
	ObjectType business.RankObjectType
	DateType   business.RankDateType
	TimePoint  time.Time
	Num        int
}
type GetActualUserCarbonRankParam struct {
	StartTime time.Time
	EndTime   time.Time
	CompanyId int
	Limit     int
	Offset    int
}
type GetActualDepartmentCarbonRankParam struct {
	StartTime time.Time
	EndTime   time.Time
	CompanyId int
	Limit     int
	Offset    int
}
type UserRankInfo struct {
	User    business.User   `json:"user"`
	IsLike  bool            `json:"isLike"`
	LikeNum int             `json:"likeNum"`
	Value   decimal.Decimal `json:"value"`
}
type GetUserRankListParam struct {
	UserId    int64
	DateType  business.RankDateType
	CompanyId int
	Limit     int
	Offset    int
}
type GetMyRankParam struct {
	UserId    int64
	DateType  business.RankDateType
	CompanyId int
}
