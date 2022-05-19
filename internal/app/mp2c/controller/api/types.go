package api

import (
	"mio/internal/app/mp2c/controller"
	"time"
)

type GetWeappQrCodeFrom struct {
	TopicId int `json:"topicId" form:"topicId" binding:"required" alias:"文章id"`
}
type GetTopicPageListForm struct {
	ID         int64 `json:"id" form:"id" binding:"gte=0" alias:"topic id"`
	TopicTagId int   `json:"topicTagId" form:"topicTagId" binding:"gte=0" alias:"标签id"`
	controller.PageFrom
}
type ChangeTopicLikeForm struct {
	TopicId int `json:"topicId" form:"topicId" binding:"required" alias:"文章id"`
}
type GetTagForm struct {
	ID int `json:"id" form:"id" binding:"gte=0" alias:"tag id"`
	controller.PageFrom
}
type GetYZMForm struct {
	Mobile string `json:"mobile" form:"mobile" binding:"required" alias:"手机号码"`
	Code   string `json:"code" form:"code"  alias:"验证码"`
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
	Ch         string  `json:"ch" form:"ch" binding:"required" alias:"渠道参数"`
	Mobile     string  `json:"mobile" form:"mobile" binding:"required" alias:"手机号码"`
	OutTradeNo string  `json:"outTradeNo" form:"outTradeNo" binding:"required" alias:"唯一订单号"`
	TotalPower float32 `json:"totalPower" form:"totalPower" binding:"required" alias:"总电量"`
	Sign       string  `json:"sign" form:"sign" binding:"required" alias:"签名"`
}

type DuibaAutoLoginForm struct {
	Path string `json:"path" form:"path"`
}
type BindMobileByCodeForm struct {
	Code string `json:"code" form:"code" binding:"required" alias:"code"`
}
type GetPointTransactionListForm struct {
	StartTime time.Time `json:"startTime" form:"startTime"  time_format:"2006-01-02 15:04:05" time_utc:"false" time_location:"Asia/Shanghai"`
	EndTime   time.Time `json:"endTime" form:"endTime"  time_format:"2006-01-02 15:04:05" time_utc:"false" time_location:"Asia/Shanghai"`
}
type UpdateStepTotalForm struct {
	EncryptedData string `json:"encryptedData" form:"encryptedData" binding:"required" alias:"encryptedData"`
	IV            string `json:"iv" form:"iv" binding:"required" alias:"iv"`
}

type AnswerQuizQuestionForm struct {
	QuestionId string `json:"questionId" form:"questionId" binding:"required" alias:"questionId"`
	Choice     string `json:"choice" form:"choice" binding:"required" alias:"choice"`
}
