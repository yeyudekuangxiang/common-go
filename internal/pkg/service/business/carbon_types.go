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
type DepartmentRankInfo struct {
	Department business.Department
	IsLike     bool            `json:"isLike"`
	LikeNum    int             `json:"likeNum"`
	Rank       int             `json:"rank"`
	Value      decimal.Decimal `json:"value"`
}
type GetUserRankListParam struct {
	UserId    int64 //用于判断对排行榜是否点赞
	DateType  business.RankDateType
	CompanyId int
	Limit     int
	Offset    int
}
type GetDepartmentRankListParam struct {
	UserId    int64 //用于判断对排行榜是否点赞
	DateType  business.RankDateType
	CompanyId int
	Limit     int
	Offset    int
}
type FindUserRankParam struct {
	UserId   int64
	DateType business.RankDateType
}
type FindDepartmentRankParam struct {
	UserId       int64 //用于判断当前用户是否对部门点赞
	DepartmentId int
	DateType     business.RankDateType
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
	Info          business.CarbonTypeInfo
	Type          business.CarbonType
	TransactionId string
}
type CompanySceneSetting struct {
	PointSetting business.PointSetting
	MaxCount     int
}

type SendCarbonCreditEvCarParam struct {
	UserId        int64
	Electricity   float64
	TransactionId string
}
type CarbonCreditEvCarResult struct {
	Credit decimal.Decimal
}
type SendCarbonCreditOnlineMeetingParam struct {
	UserId           int64
	OneCityDuration  time.Duration
	ManyCityDuration time.Duration
	TransactionId    string
}
type SendCarbonCreditOnlineMeetingResult struct {
	OneCityCredit  decimal.Decimal
	ManyCityCredit decimal.Decimal
}

type SendCarbonCreditSaveWaterElectricityParam struct {
	UserId        int64
	Water         int64
	Electricity   int64
	TransactionId string
}
type SendCarbonCreditSaveWaterElectricityResult struct {
	WaterCredit       decimal.Decimal
	ElectricityCredit decimal.Decimal
}
type SendCarbonCreditSavePublicTransportParam struct {
	UserId        int64
	Bus           int64
	Metro         int64
	TransactionId string
}
type SendCarbonCreditSavePublicTransportResult struct {
	BusCredits   decimal.Decimal
	MetroCredits decimal.Decimal
}
type CarbonResult struct {
	Credit decimal.Decimal
	Point  int
}
