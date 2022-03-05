package service

import (
	"mio/model/entity"
)

type GetTopicistParam struct {
	ID         int `json:"id" form:"id"`
	TopicTagId int `json:"topicTagId" form:"topicTagId"`
	Offset     int `json:"offset"`
	Limit      int `json:"limit"` //limit为0时不限制数量
}
type TopicDetail struct {
	entity.Topic
	IsLike bool `json:"isLike"`
}