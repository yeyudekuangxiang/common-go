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
	UserId   int64
}

type GetTopicPageListBy struct {
	ID         int64              `json:"id"`
	TopicTagId int                `json:"topicTagId"`
	Offset     int                `json:"offset"`
	Limit      int                `json:"limit"`  //limit为0时不限制数量
	UserId     int64              `json:"userId"` // 用于查询用户对帖子是否点赞
	OrderBy    entity.OrderByList `json:"orderBy"`
}

type GetTagPageListBy struct {
	ID      int                `json:"id"`
	Offset  int                `json:"offset"`
	Limit   int                `json:"limit"` //limit为0时不限制数量
	OrderBy entity.OrderByList `json:"orderBy"`
}

type GetUserListBy struct {
	Mobile  string
	Mobiles []string
	Source  entity.UserSource
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
type GetTopicListBy struct {
	TopicIds []int64
}
type FindTopicFlowBy struct {
	TopicId int64
	UserId  int64
}
type GetTopicFlowPageListBy struct {
	Offset     int
	Limit      int
	UserId     int64
	TopicId    int64
	TopicTagId int
}
type GetProductItemListBy struct {
	ItemIds []string
}
type CheckStockItem struct {
	ItemId string
	Count  int
}
