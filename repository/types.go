package repository

import (
	"mio/model"
	"mio/model/entity"
)

type GetUserBy struct {
	OpenId string
	Source entity.UserSource
	Mobile string
}
type FindTopicLikeBy struct {
	TopicId int
	UserId  int
}
type GetTopicLikeListBy struct {
	TopicIds []int64
	UserId   int
}

type GetTopicPageListBy struct {
	ID         int `json:"id"`
	TopicTagId int `json:"topicTagId"`
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`  //limit为0时不限制数量
	UserId     int `json:"userId"` // 用于查询用户对帖子是否点赞
}

type GetTagPageListBy struct {
	ID     int `json:"id"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"` //limit为0时不限制数量
}

type GetUserListBy struct {
	Mobile  string
	UserIds []int64
}

type FindPointBy struct {
	OpenId string
}
type GetPointTransactionListBy struct {
	OpenId string
}
type FindPointTransactionCountLimitBy struct {
	OpenId          string
	TransactionType entity.PointTransactionType
	TransactionDate model.Date
}
