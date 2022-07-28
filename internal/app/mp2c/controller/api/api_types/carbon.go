package api_types

import "mio/internal/pkg/model/entity"

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
	Nickname  string
	AvatarUrl string
	Uid       int64
}

type CarbonTransactionMyBank struct {
	User    CarbonUser
	Carbon  float64
	Rank    string
	OverPer string
}

type CarbonTransactionBank struct {
	User     CarbonUser
	Carbon   float64
	Rank     int64
	IsFriend bool
}

//form

type GetCarbonTransactionBankForm struct {
	Scene string `json:"scene" form:"scene" binding:"oneof=home event topic" alias:"banner场景"`
}
type GetCarbonTransactionMyBankForm struct {
	Scene string `json:"scene" form:"scene" binding:"oneof=home event topic" alias:"banner场景"`
}
type GetCarbonTransactionInfoForm struct {
	Scene string `json:"scene" form:"scene" binding:"oneof=home event topic" alias:"banner场景"`
}
type GetCarbonTransactionClassifyForm struct {
	Scene string `json:"scene" form:"scene" binding:"oneof=home event topic" alias:"banner场景"`
}
type GetCarbonTransactionHistoryForm struct {
	Scene string `json:"scene" form:"scene" binding:"oneof=home event topic" alias:"banner场景"`
}

// DTO

type CreateCarbonTransactionDto struct {
	OpenId  string `binding:"required"`
	UserId  int64
	Type    entity.CarbonTransactionType
	Value   float64
	Info    string
	AdminId int
}

type GetCarbonTransactionBankDto struct {
	Offset int64
	Limit  int64
	UserId int64
}
type GetCarbonTransactionMyBankDto struct {
	UserId int64
}
type GetCarbonTransactionInfoDto struct {
	UserId int64
}
type GetCarbonTransactionClassifyDto struct {
	UserId int64
}
type GetCarbonTransactionHistoryDto struct {
	UserId int64
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
