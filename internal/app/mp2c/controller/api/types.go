package api

import (
	"mio/internal/app/mp2c/controller"
	"time"
)

type GetWeappQrCodeFrom struct {
	TopicId int64 `json:"topicId" form:"topicId" binding:"required" alias:"文章id"`
}
type GetTopicPageListForm struct {
	ID         int64 `json:"id" form:"id" binding:"gte=0" alias:"topic id"`
	TopicTagId int64 `json:"topicTagId" form:"topicTagId" binding:"gte=0" alias:"标签id"`
	controller.PageFrom
}
type ChangeTopicLikeForm struct {
	TopicId int `json:"topicId" form:"topicId" binding:"required" alias:"文章id"`
}
type GetTagForm struct {
	ID int64 `json:"id" form:"id" binding:"gte=0" alias:"tag id"`
	controller.PageFrom
}
type GetYZMForm struct {
	Mobile string `json:"mobile" form:"mobile" binding:"required" alias:"手机号码"`
	Code   string `json:"code" form:"code"  alias:"验证码"`
	Cid    int64  `json:"cid" form:"cid"  alias:"渠道来源"`
}

type CreateQrcodeForm struct {
	Src string `json:"src" form:"src" binding:"required" alias:"跳转链接"`
}
type SubmitOrderForGreenForm struct {
	AddressId string `json:"addressId" form:"addressId" binding:"required" alias:"地址"`
}
type GetOCRForm struct {
	Src string `json:"src" form:"src" binding:"required" alias:"图片地址"`
}
type GetChargeForm struct {
	Ch         string `json:"ch" form:"ch" binding:"required" alias:"渠道参数"`
	Mobile     string `json:"mobile" form:"mobile" binding:"required" alias:"手机号码"`
	OutTradeNo string `json:"outTradeNo" form:"outTradeNo" binding:"required" alias:"唯一订单号"`
	TotalPower string `json:"totalPower" form:"totalPower" binding:"required" alias:"统计总量"`
	Sign       string `json:"sign" form:"sign" binding:"required" alias:"签名"`
}
type ChangeChargeExceptionForm struct {
	Ch string `json:"ch" form:"ch" binding:"required" alias:"渠道参数"`
}

type RecyclePushForm struct {
	Type                string `json:"type" form:"type" binding:"required"`                     //业务类型 1：回首订单成功
	OrderNo             string `json:"orderNo" form:"orderNo" binding:"required"`               //订单号，同类型同订单视为重复订单
	Name                string `json:"name,omitempty" form:"name" binding:"required"`           //type = 1，回收物品名称
	OolaUserId          string `json:"oolaUserId" form:"oolaUserId" binding:"required"`         //噢啦平台用户id
	ClientId            string `json:"clientId" form:"clientId" binding:"required"`             //lvmiao用户id
	CreateTime          string `json:"createTime" form:"createTime" binding:"required"`         //订单创建时间
	CompletionTime      string `json:"completionTime" form:"completionTime" binding:"required"` //订单完成时间
	BeanNum             string `json:"beanNum" form:"beanNum"`
	Sign                string `json:"sign" form:"sign" binding:"required"`                      //加密串
	ProductCategoryName string `json:"productCategoryName,omitempty" form:"productCategoryName"` //物品所属分类名称
	Qua                 string `json:"qua,omitempty" form:"qua"`                                 //用户下单时的数量&重量
	Unit                string `json:"unit,omitempty" form:"unit"`                               //与下单数量&重量关联的计量单位 如：公斤，个 等
}

type QueryTokenForm struct {
	OperatorID     string `json:"operatorID" form:"operatorID" alias:"运营商标识" binding:"required"`
	OperatorSecret string `json:"operatorSecret" form:"operatorSecret" alias:"运营商密钥" binding:"required"`
}

type DuibaAutoLoginForm struct {
	Path string `json:"path" form:"path"`
}
type BindMobileByCodeForm struct {
	Code      string `json:"code" form:"code" binding:"required" alias:"code"`
	InvitedBy string `json:"invitedBy" form:"invitedBy" alias:"invitedBy"`
}
type GetPointTransactionListForm struct {
	StartTime time.Time `json:"startTime" form:"startTime"  time_format:"2006-01-02" time_utc:"false" time_location:"Asia/Shanghai"`
	EndTime   time.Time `json:"endTime" form:"endTime"  time_format:"2006-01-02" time_utc:"false" time_location:"Asia/Shanghai"`
}
type UpdateStepTotalForm struct {
	EncryptedData string `json:"encryptedData" form:"encryptedData" binding:"required" alias:"encryptedData"`
	IV            string `json:"iv" form:"iv" binding:"required" alias:"iv"`
}

type AnswerQuizQuestionForm struct {
	QuestionId string `json:"questionId" form:"questionId" binding:"required" alias:"questionId"`
	Choice     string `json:"choice" form:"choice" binding:"required" alias:"choice"`
}
type UploadPointCollectImageForm struct {
	PointCollectType string `json:"pointCollectType" form:"pointCollectType" binding:"oneof=COFFEE_CUP BIKE_RIDE DIDI POWER_REPLACE" alias:"类型"`
}
type PointCollectForm struct {
	ImgUrl           string `json:"imgUrl" form:"imgUrl" binding:"required" alias:"图片"`
	PointCollectType string `json:"pointCollectType" form:"pointCollectType" binding:"oneof=COFFEE_CUP BIKE_RIDE DIDI POWER_REPLACE" alias:"类型"`
}

type NewCollectForm struct {
	ImgUrl           string `json:"imgUrl" form:"imgUrl" binding:"required" alias:"图片"`
	PointCollectType string `json:"pointCollectType" form:"pointCollectType" alias:"类型"`
}

type CollectType struct {
	PointCollectType string `json:"pointCollectType" form:"pointCollectType" binding:"required" alias:"类型"`
}

type UpdateUserInfoForm struct {
	Nickname    string  `json:"nickname" form:"nickname"`
	Avatar      string  `json:"avatar" form:"avatar"`
	Gender      *string `json:"gender" form:"gender" binding:"omitempty,oneof=MALE FEMALE"`
	Birthday    *string `json:"birthday" form:"birthday" binding:"omitempty,datetime=2006-01-02"`
	PhoneNumber *string `json:"phoneNumber" form:"phoneNumber" binding:"omitempty,gt=0"`
}
type DuiBaNoLoginH5Form struct {
	ActivityId string `json:"activityId" form:"activityId" `
}

type CreateTopicForm struct {
	Title   string   `json:"title" form:"title" alias:"title" binding:"required,min=2,max=64"`
	Content string   `json:"content" form:"content" alias:"content" binding:"min=0,max=10000"`
	Images  []string `json:"images" form:"images" alias:"images" binding:"required,min=1,max=12"`
	TagIds  []int64  `json:"tagIds" form:"tagIds" alias:"tagIds" binding:"min=0,max=2"`
}

type UpdateTopicForm struct {
	ID int64 `json:"id" form:"id" alias:"id" binding:"required,gte=1"`
	CreateTopicForm
}

type IdForm struct {
	ID int64 `json:"id" form:"id" alias:"id" binding:"required,gte=1"`
}

type ListFormById struct {
	ID int64 `json:"id" form:"id" alias:"id" binding:"required,gte=1"`
	controller.PageFrom
}

type ListFormByLastId struct {
	ID       int64 `json:"id" form:"id" alias:"id" binding:"required,gte=1"`
	LastId   int64 `json:"lastId" form:"lastId" alias:"lastId" binding:"required,gte=1"`
	PageSize int   `json:"pageSize" form:"pageSize" binding:"gt=0" alias:"每页数量"`
}

type CommentCreateForm struct {
	Message string `json:"message" form:"message" alias:"message" binding:"required,min=1"`
	Root    int64  `json:"root" form:"root" alias:"root" binding:"min=0"`
	Parent  int64  `json:"parent" form:"parent" alias:"parent" binding:"min=0"`
	ObjId   int64  `json:"objId" form:"objId" alias:"objId" binding:"required,min=1"`
}

type CommentEditForm struct {
	CommentId int64  `json:"commentId" form:"commentId" alias:"commentId" binding:"required,min=1"`
	Message   string `json:"message" form:"message" alias:"message" binding:"required,min=1"`
}

type ChangeCommentLikeForm struct {
	CommentId int64 `json:"commentId" form:"commentId" binding:"required" alias:"评论id"`
}
