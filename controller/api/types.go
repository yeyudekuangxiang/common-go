package api

import "mio/controller"

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
