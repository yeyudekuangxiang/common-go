package activity

import (
	"mio/internal/app/mp2c/controller"
)

type GetBocRecordListForm struct {
	ApplyStatus int8 `json:"applyStatus" form:"applyStatus" binding:"oneof=0 1 2 3 4" alias:"参与状态"`
	controller.PageFrom
}
type AddBocRecordFrom struct {
	ShareUserId int64  `json:"shareUserId" form:"shareUserId" binding:"gte=0" alias:"分享者ID"`
	Source      string `json:"source" form:"source" alias:"用户来源"`
}

type AnswerBocQuestionFrom struct {
	Right int8 `json:"right" form:"right" binding:"oneof=2 3" alias:"答题结果"`
}

type ApplySendBonusForm struct {
	Type string `json:"type" form:"type" binding:"oneof=apply bind boc" alias:"奖励类型"`
}
type ReportInvitationRecordForm struct {
	UserId int64 `json:"userId" form:"userId"`
}
type ExchangeGiftForm struct {
	AddressId string `json:"addressId" form:"addressId" binding:"required" alias:"地址"`
}

type GMAnswerQuestion struct {
	Title   string `json:"title" form:"title" binding:"required" alias:"问题标题"`
	IsRight bool   `json:"isRight" form:"isRight"`
	Answer  string `json:"answer" form:"answer" binding:"required" alias:"答案"`
}
type ZeroAutoLoginForm struct {
	Short string `json:"short" form:"short"`
}
type ZeroStoreUrlForm struct {
	Url string `json:"Url" form:"Url" binding:"required" alias:"url"`
}
type DuiBaAutoLoginForm struct {
	ActivityId string `json:"activityId" form:"activityId" `
	Short      string `json:"short" form:"short"`
}
type DuiBaStoreUrlForm struct {
	ActivityId string `json:"activityId" form:"activityId" `
	Url        string `json:"Url" form:"Url" binding:"required" alias:"url"`
}
type GetDuiBaActivityQrForm struct {
	ActivityId string `json:"activityId" form:"activityId" binding:"required" alias:"activityId"`
	Password   string `json:"password" form:"password" binding:"required" alias:"password"`
}
