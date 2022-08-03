package repository

import (
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"time"
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
	TopicTagId int64              `json:"topicTagId"`
	Offset     int                `json:"offset"`
	Status     int                `json:"status"` //0全部 1待审核 2审核失败 3已发布 4已下架
	IsTop      int                `json:"isTop"`
	IsEssence  int                `json:"isEssence"`
	Limit      int                `json:"limit"`  //limit为0时不限制数量
	UserId     int64              `json:"userId"` // 用于查询用户对帖子是否点赞
	OrderBy    entity.OrderByList `json:"orderBy"`
}

type TopicListRequest struct {
	ID        int64  `json:"id"` //帖子id
	Title     string `json:"title"`
	UserId    int64  `json:"userId"`
	UserName  string `json:"userName"`
	TagId     int64  `json:"tagId"`
	Status    int    `json:"status"`    //0全部 1待审核 2审核失败 3已发布 4已下架
	IsTop     int    `json:"isTop"`     //是否置顶
	IsEssence int    `json:"isEssence"` //是否
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
}

type GetTagPageListBy struct {
	//	ID          int                `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Offset      int                `json:"offset"`
	Limit       int                `json:"limit"` //limit为0时不限制数量
	OrderBy     entity.OrderByList `json:"orderBy"`
}

type CreateTag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}
type UpdateTag struct {
	ID int64 `json:"id"`
	CreateTag
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
	StartTime  time.Time
	EndTime    time.Time
	Risk       int
}

type GetUserPageListBy struct {
	User    GetUserListBy
	Offset  int    `json:"offset"`
	Limit   int    `json:"limit"` //limit为0时不限制数量
	OrderBy string `json:"orderBy"`
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
	AdminId   int
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

type FindTopicFlowBy struct {
	TopicId int64
	UserId  int64
}
type GetTopicFlowPageListBy struct {
	Offset     int
	Limit      int
	UserId     int64
	TopicId    int64
	TopicTagId int64
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
	OpenId        string
	Day           model.Time
	RecordedEpoch int64
	OrderBy       entity.OrderByList
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
type GetStepHistoryListBy struct {
	OpenId            string
	RecordedEpochs    []int64
	StartRecordedTime model.Time // >=
	EndRecordedTime   model.Time //<=
	OrderBy           entity.OrderByList
}
type GetStepHistoryPageListBy struct {
	GetStepHistoryListBy
	Limit  int
	Offset int
}
type FindCouponBy struct {
	CouponTypeId string
	CouponId     string
}
type FindCouponTypeBy struct {
	CouponTypeId string
}
type FindDuiBaPointAddLogBy struct {
	OrderNum string
}
type FindCertificateBy struct {
	ProductItemId string
	CertificateId string
}

type FindUserChannelBy struct {
	Cid  int64
	Code string
}

type GetUserChannelPageListBy struct {
	Name    string             `json:"name" `
	Code    string             `json:"code" `
	Cid     int                `json:"cid"`
	Pid     int                `json:"pid"`
	Offset  int                `json:"offset"`
	Limit   int                `json:"limit"` //limit为0时不限制数量
	OrderBy entity.OrderByList `json:"orderBy"`
}

type FindTopicBy struct {
	TopicIds []int64
	Status   int `json:"status"` //0全部 1待审核 2审核失败 3已发布 4已下架
}
