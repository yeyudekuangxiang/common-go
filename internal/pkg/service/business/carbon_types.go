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
	User    business.ShortUser `json:"user"`
	IsLike  bool               `json:"isLike"`
	LikeNum int                `json:"likeNum"`
	Rank    int                `json:"rank"`
	Value   decimal.Decimal    `json:"value"`
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
type ChangeLikeStatusParam struct {
	Pid        int64
	ObjectType business.RankObjectType
	DateType   business.RankDateType
	UserId     int64
}
type FindCompanyCarbonSceneParam struct {
	CompanyId     int
	CarbonSceneId int
}
type FindCarbonCreditsLimitLogParam struct {
	TimePoint time.Time
	Type      business.CarbonType
	UserId    int64
}
type UpdateOrCreateCarbonCreditsLimitLogParam struct {
	Type            business.CarbonType
	UserId          int64
	AddCurrentValue decimal.Decimal //增加的积分数量
	TimePoint       time.Time
}
type createOrUpdateCarbonCreditParam struct {
	UserId    int64
	AddCredit decimal.Decimal
}
type SendCarbonCreditParam struct {
	UserId        int64
	AddCredit     decimal.Decimal
	Info          string
	Type          business.CarbonType
	TransactionId string
}

//发放碳积分
type SendCarbonCreditEvCarParam struct {
}
