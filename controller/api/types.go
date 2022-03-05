package api

import "mio/controller"

type GetWeappQrCodeFrom struct {
	TopicId int `json:"topicId" form:"topicId" binding:"required" alias:"文章id"`
}
type GetTopicPageListForm struct {
	ID         int `json:"id" form:"id" binding:"gte=0" alias:"topic id"`
	TopicTagId int `json:"topicTagId" form:"topicTagId" binding:"gte=0" alias:"标签id"`
	controller.PageFrom
}
type ChangeTopicLikeForm struct {
	TopicId int `json:"topicId" form:"topicId" binding:"required" alias:"文章id"`
}
