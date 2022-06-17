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
	ThirdParty string `json:"thirdParty" form:"thirdParty"`
}
type DuiBaStoreUrlForm struct {
	ActivityId string `json:"activityId" form:"activityId" `
	Url        string `json:"Url" form:"Url" binding:"required" alias:"url"`
}
type GetDuiBaActivityQrForm struct {
	ActivityId string `json:"activityId" form:"activityId" binding:"required" alias:"activityId"`
	Password   string `json:"password" form:"password" binding:"required" alias:"password"`
}

type GDDbActivityHomePageForm struct {
	//UserId   int64 `json:"userId" form:"userId" alias:"userId"`
	InviteId int64 `json:"inviteId" form:"inviteId" alias:"inviteId"`
}

type GDDbActivitySchoolForm struct {
	UserName string `json:"userName" form:"userName" binding:"required" alias:"userName"`
	SchoolId int64  `json:"schoolId" form:"schoolId" binding:"required" alias:"schoolId"`
	//ProvinceId  int64  `json:"provinceId" form:"provinceId" binding:"required" alias:"provinceId"`                   //省id
	CityId int64 `json:"cityId" form:"cityId" binding:"required" alias:"cityId"` //市id
	//AreaId      int64  `json:"areaId" form:"areaId" binding:"required" alias:"areaId"`                               //区id
	GradeId     int64  `json:"gradeId" form:"gradeId" binding:"required" alias:"gradeId"`                            //年级id
	ClassNumber uint32 `json:"classNumber" form:"classNumber" binding:"required,max=1000,min=1" alias:"classNumber"` //班级号码
}

type GDDbSelectSchoolForm struct {
	SchoolName string `json:"schoolName" form:"schoolName" alias:"schoolName"`
	CityId     int64  `json:"cityId" form:"cityId" binding:"gte=0" alias:"cityId"`
	GradeId    int64  `json:"gradeId" form:"gradeId" binding:"gte=0" alias:"gradeId"`
}

type GDDbCreateSchoolForm struct {
	SchoolName string `json:"schoolName" form:"schoolName" binding:"required" alias:"schoolName"`
	CityId     int64  `json:"cityId" form:"cityId" binding:"required,number" alias:"cityId"`
	GradeType  int    `json:"gradeType" form:"gradeType" binding:"gte=0" alias:"gradeType"`
}
