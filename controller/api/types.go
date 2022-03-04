package api

type GetWeappQrCodeFrom struct {
	OpenId  string `json:"openId" form:"openId" binding:"required" alias:"openId"`
	TopicId int    `json:"topicId" form:"topicId" binding:"required" alias:"文章id"`
}
type GetTopicListForm struct {
	ID         *int `json:"id" form:"id" binding:"omitempty,gt=0" alias:"topic id"`
	TopicTagId *int `json:"topicTagId" form:"topicTagId" binding:"omitempty,gt=0" alias:"标签id"`
}
