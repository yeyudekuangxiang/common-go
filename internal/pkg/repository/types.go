package repository

import (
	"mio/internal/pkg/model"
	entity2 "mio/internal/pkg/model/entity"
)

type GetUserBy struct {
	OpenId  string
	Source  entity2.UserSource
	Mobile  string
	UnionId string
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
	ID         int64               `json:"id"`
	TopicTagId int                 `json:"topicTagId"`
	Offset     int                 `json:"offset"`
	Status     int                 `json:"status"` //0全部 1待审核 2审核失败 3已发布 4已下架
	Limit      int                 `json:"limit"`  //limit为0时不限制数量
	UserId     int64               `json:"userId"` // 用于查询用户对帖子是否点赞
	OrderBy    entity2.OrderByList `json:"orderBy"`
}

type GetTagPageListBy struct {
	ID      int                 `json:"id"`
	Offset  int                 `json:"offset"`
	Limit   int                 `json:"limit"` //limit为0时不限制数量
	OrderBy entity2.OrderByList `json:"orderBy"`
}

type GetUserListBy struct {
	Mobile  string
	Mobiles []string
	Source  entity2.UserSource
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
	TransactionType entity2.PointTransactionType
	TransactionDate model.Date
}
type GetTopicListBy struct {
	TopicIds []int64
	Status   int `json:"status"` //0全部 1待审核 2审核失败 3已发布 4已下架
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
	Status     entity2.TopicStatus `json:"status"` //直接传0值表示全部
}
type GetProductItemListBy struct {
	ItemIds []string
}
type CheckStockItem struct {
	ItemId string
	Count  int
}
type FindPointTransactionBy struct {
	TransactionId string
}
