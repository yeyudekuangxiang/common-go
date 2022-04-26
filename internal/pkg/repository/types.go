package repository

import (
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
)

type GetUserBy struct {
	OpenId     string
	Source     entity.UserSource
	Mobile     string //手机号精确匹配
	LikeMobile string //手机号模糊匹配
	UnionId    string
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
	Status     int                `json:"status"` //0全部 1待审核 2审核失败 3已发布 4已下架
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
	Mobile     string
	Mobiles    []string
	Source     entity.UserSource
	UserIds    []int64
	Nickname   string //模糊查询
	LikeMobile string //手机号模糊查询
	UserId     int64
	OpenId     string
}

type FindPointBy struct {
	OpenId string
}
type GetPointTransactionListBy struct {
	OpenId    string
	StartTime model.Time
	EndTime   model.Time
	OrderBy   entity.OrderByList
	Type      entity.PointTransactionType
}
type GetPointTransactionPageListBy struct {
	OpenIds   []string
	StartTime model.Time
	EndTime   model.Time
	OrderBy   entity.OrderByList
	Type      entity.PointTransactionType
	Types     []entity.PointTransactionType
	Offset    int
	Limit     int
}
type FindPointTransactionCountLimitBy struct {
	OpenId          string
	TransactionType entity.PointTransactionType
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
	Status     entity.TopicStatus `json:"status"` //直接传0值表示全部
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
type FindStepHistoryBy struct {
	OpenId  string
	Day     model.Time
	OrderBy entity.OrderByList
}
type FindStepBy struct {
	OpenId string
}
type GetStepListBy struct {
	OpenId       string
	RecordedTime model.Time
}
type GetFileExportPageListBy struct {
	Type           entity.FileExportType
	AdminId        int64
	Status         entity.FileExportStatus
	StartCreatedAt model.Time
	EndCreatedAt   model.Time
	OrderBy        entity.OrderByList
	Offset         int
	Limit          int
}
type GetAdminListBy struct {
}
type FindAdminBy struct {
	Account string
}

type FindOaAuthWhiteBy struct {
	Domain string
	AppId  string
}
