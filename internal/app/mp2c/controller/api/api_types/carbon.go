package api_types

import (
	"mio/internal/app/mp2c/controller"
	"mio/internal/pkg/model/entity"
)

type PageFrom struct {
	Page     int `json:"page" form:"page" binding:"gt=0" alias:"页码"`
	PageSize int `json:"pageSize" form:"pageSize" binding:"gt=0" alias:"每页数量"`
}

func (p PageFrom) Limit() int {
	return p.PageSize
}

func (p PageFrom) Offset() int {
	return (p.Page - 1) * p.PageSize
}

//service返回的结构

type RedeemCouponWithTransactionResult struct {
}

type CarbonUser struct {
	Nickname  string `json:"nickname"`
	AvatarUrl string `json:"avatarUrl"`
	Uid       int64  `json:"uid"`
}

type CarbonTransactionMyBank struct {
	User    CarbonUser `json:"user"`
	Carbon  string     `json:"carbon"`
	Rank    string     `json:"rank"`
	OverPer string     `json:"overPer"`
}

type CarbonTransactionBank struct {
	User     CarbonUser `json:"user"`
	Carbon   string     `json:"carbon"`
	Rank     int64      `json:"rank"`
	IsFriend bool       `json:"isFriend"`
}

type CarbonTransactionInfo struct {
	User           CarbonUser `json:"user"`
	Carbon         string     `json:"carbon"`
	CarbonToday    string     `json:"carbonToday"`
	RegisterDayNum int        `json:"registerDayNum"`
	TreeNum        string     `json:"treeNum"`
	TreeNumMsg     string     `json:"treeNumMsg"`
}

type CarbonTransactionClassify struct {
	List  []CarbonTransactionClassifyList
	Cover string  `json:"cover"`
	Total float64 `json:"total"`
}

type CarbonTransactionClassifyList struct {
	Key string  `json:"key"`
	Val float64 `json:"val"`
}

//form

type GetCarbonTransactionBankForm struct {
	controller.PageFrom
}

type GetCarbonTransactionCreateForm struct {
	SerialNumber string `json:"serialNumber" form:"serialNumber"  binding:"required" alias:"流水号"`
	CarbonType   string `json:"carbonType" form:"carbonType" binding:"required" alias:"类型"`
	CarbonValue  string `json:"carbonValue" form:"carbonValue" binding:"required" alias:"碳量"`
	Uid          string `json:"uid" form:"uid" binding:"required" alias:"uid"`
	Time         string `json:"time" form:"time" binding:"datetime=2006-01-02 15:04:05" alias:"时间"`
	Sign         string `json:"sign" form:"sign" binding:"required" alias:"sign"`
}

// DTO

type CreateCarbonTransactionDto struct {
	OpenId   string `binding:"required"`
	UserId   int64
	Type     entity.CarbonTransactionType
	Value    float64
	Info     string
	AdminId  int
	Ip       string
	AddValue float64
}

type GetCarbonTransactionBankDto struct {
	Offset int
	Limit  int
	UserId int64
}
type GetCarbonTransactionMyBankDto struct {
	UserId int64
}
type GetCarbonTransactionInfoDto struct {
	UserId int64
}
type GetCarbonTransactionClassifyDto struct {
	UserId    int64
	StartTime string
	EndTime   string
}
type GetCarbonTransactionHistoryDto struct {
	UserId    int64
	StartTime string
	EndTime   string
}

//VO

type CarbonBankVo struct {
	MyCarbon   CarbonTransactionMyBank
	CarbonUser CarbonUser
}

type BankCarbonTransactionParam struct {
	OpenId        string `binding:"required"`
	UserId        int64
	Type          entity.CarbonTransactionType `binding:"required"`
	Value         float64
	Info          string
	AdminId       int
	Offset, Count int64
}

type CarbonTransactionClassifyVO struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

type CarbonTransactionHistoryVO struct {
	VDate    string  `json:"vDate"`
	Value    float64 `json:"value"`
	ValueStr string  `json:"valueStr"`
}

type CarbonTransactionHistoryDateVO struct {
	VDate string `json:"vDate"`
}

type CarbonTransactionHistoryValueVO struct {
	Value    float64 `json:"value"`
	ValueStr string  `json:"valueStr"`
}
