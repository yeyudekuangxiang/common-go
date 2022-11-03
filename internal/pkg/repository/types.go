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

type UpdateUserRisk struct {
	UserIdSlice []string `json:"UserIdSlice"`
	OpenIdSlice []string `json:"OpenIdSlice"`
	PhoneSlice  []string `json:"PhoneSlice"`
	Risk        int      `json:"Risk"`
}

type FindTopicLikeBy struct {
	TopicId int64
	UserId  int64
}
type GetTopicLikeListBy struct {
	TopicIds []int64
	UserIds  []int64
	UserId   int64
	TopicId  int64
	Status   int
}

type GetTopicPageListBy struct {
	ID          int64              `json:"id"`
	Ids         []int64            `json:"ids"`
	TopicTagId  int64              `json:"topicTagId"`
	Offset      int                `json:"offset"`
	Status      int                `json:"status"` //0全部 1待审核 2审核失败 3已发布 4已下架
	IsTop       int                `json:"isTop"`
	IsEssence   int                `json:"isEssence"`
	Limit       int                `json:"limit"`  // limit为0时不限制数量
	UserId      int64              `json:"userId"` // 用于查询用户对帖子是否点赞
	OrderByList entity.OrderByList `json:"orderByList"`
	OrderBy     entity.OrderBy     `json:"orderBy"`
	Order       string             `json:"order"`
}

type GetTopicCountBy struct {
	TopicTagId int64  `json:"topicTagId"`
	Status     int    `json:"status"` //0全部 1待审核 2审核失败 3已发布 4已下架
	IsTop      int    `json:"isTop"`
	IsEssence  int    `json:"isEssence"`
	UserId     int64  `json:"userId"` // 用于查询用户对帖子是否点赞
	OpenId     string `json:"openId"`
}

type TopicListRequest struct {
	ID         int64    `json:"id"` //帖子id
	Title      string   `json:"title"`
	UserId     int64    `json:"userId"`
	UserName   string   `json:"userName"`
	IsPartners int      `json:"isPartners"`
	Position   string   `json:"position"`
	TagId      int64    `json:"tagId"`
	TagIds     []string `json:"tagIds"`
	Status     int      `json:"status"`    //0全部 1待审核 2审核失败 3已发布 4已下架
	IsTop      int      `json:"isTop"`     //是否置顶
	IsEssence  int      `json:"isEssence"` //是否
	Offset     int      `json:"offset"`
	Limit      int      `json:"limit"`
}

type GetTagPageListBy struct {
	ID          int64              `json:"id"`
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
	Mobile     string              `json:"mobile,omitempty"`
	Mobiles    []string            `json:"mobiles,omitempty"`
	Source     entity.UserSource   `json:"source,omitempty"`
	UserIds    []int64             `json:"userIds,omitempty"`
	Nickname   string              `json:"nickname,omitempty"`   //模糊查询
	LikeMobile string              `json:"likeMobile,omitempty"` //手机号模糊查询
	UserId     int64               `json:"userId,omitempty"`
	Status     int                 `json:"status,omitempty"` //0全部 1正常 2禁言 3封号
	OpenId     string              `json:"openId,omitempty"`
	StartTime  time.Time           `json:"startTime"`
	EndTime    time.Time           `json:"endTime"`
	Risk       int                 `json:"risk,omitempty"`
	Position   entity.UserPosition `json:"position"`
	Partners   entity.Partner      `json:"partner"` //0:全部 1:乐活家 2:非乐活家
	Auth       int                 `json:"auth"`
}

type GetUserIdentifyInfoBy struct {
	ChannelName     string     `json:"user_channel"`
	CityName        string     `json:"city_name"`
	Balance         string     `json:"balance"`
	InvitedByOpenid string     `json:"invited_by_openid"`
	Openid          string     `json:"openid"`
	NickName        string     `json:"nick_name"`
	ChannelTypeName string     `json:"channel_type_name"`
	Time            model.Time `json:"time"`
	Source          string     `json:"source"`
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

type FindListPoint struct {
	OpenIds []string
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
type GetPointTransactionCountBy struct {
	AdminId int
	OpenId  string
	OpenIds []string
	Type    entity.PointTransactionType
	Types   []entity.PointTransactionType
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
	Status     entity.TopicStatus `json:"status"` // 直接传'0'值表示全部
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
	OpenId        string
	Type          string
	Note          string
	AdditionInfo  string
}
type FindPointTransactionByValue struct {
	OpenId     string
	Type       string
	ChangeType string
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
	OpenId       string
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
	Name     string             `json:"name" `
	Code     string             `json:"code" `
	Cid      int                `json:"cid"`
	CidSlice []int64            `json:"cidSlice"`
	Pid      int                `json:"pid"`
	Offset   int                `json:"offset"`
	Limit    int                `json:"limit"` //limit为0时不限制数量
	OrderBy  entity.OrderByList `json:"orderBy"`
}

type FindTopicBy struct {
	TopicIds []int64
	Status   int `json:"status"` //0全部 1待审核 2审核失败 3已发布 4已下架
}

type GetPointTransactionListByQun struct {
	OpenId    string
	StartTime string
	EndTime   string
	Types     []string
}

type FindCarbonBy struct {
	OpenId string
}

type FindCommentLikeBy struct {
	CommentId int64
	UserId    int64
}
type GetCommentLikeListBy struct {
	CommentIds []int64
	UserId     int64
}

type GetRedeemCodeBy struct {
	CodeId   string
	CouponId string
}

type GetScenePrePoint struct {
	PlatformKey    string    `json:"platformKey"`
	PlatformUserId string    `json:"platformUserId"`
	OpenId         string    `json:"openId"`
	StartTime      time.Time `json:"startTime"`
	EndTime        time.Time `json:"endTime"`
	Id             int64     `json:"id"`
	Ids            []int64   `json:"ids"`
	Status         int       `json:"status"`
}

type UpScenePrePoint struct {
	Status uint `json:"status,omitempty"`
}

type GetSceneUserOne struct {
	Id             int64     `json:"id"`
	PlatformKey    string    `json:"platformKey"`
	PlatformUserId string    `json:"platformUserId"`
	OpenId         string    `json:"openId"`
	Phone          string    `json:"phone"`
	UnionId        string    `json:"unionId"`
	StartTime      time.Time `json:"startTime"`
	EndTime        time.Time `json:"endTime"`
}

type GetSceneCallback struct {
	ID             int
	PlatformKey    string
	PlatformUserId string
	OpenId         string
	BizId          string
	SourceKey      string
	StartTime      string
	EndTime        string
}

type UpdateUserInfoParam struct {
	UserId       int64
	Nickname     string
	Avatar       string
	Gender       *entity.UserGender
	Birthday     *time.Time
	PhoneNumber  *string
	Position     string
	Status       int
	Partners     int
	Auth         int
	Introduction string
}

type SendMessage struct {
	SendId   int64  `json:"sendId"`
	RecId    int64  `json:"recId"`
	Type     int    `json:"type"`     // 1点赞 2评论 3回复 4发布 5精选 6违规 7合作社 8 商品 9订单
	TurnType int    `json:"turnType"` // 1文章 2评论 3订单 4商品
	TurnId   int64  `json:"turnId"`
	Message  string `json:"message"`
}

type GetMessage struct {
	Type int `json:"type"`
}

type FindMessageParams struct {
	MessageIds []int64   `json:"messageIds"`
	MessageId  int64     `json:"messageId"`
	SendId     int64     `json:"sendId"`
	RecId      int64     `json:"recId"`
	Type       int       `json:"type"` // 1互动消息 2酷喵圈社区 3二手交易
	Types      []int     `json:"types"`
	Status     int       `json:"status" default:"1"` // 1未读 2已读
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	Limit      int       `json:"limit"`
	Offset     int       `json:"offset"`
}

type SetHaveReadMessageParams struct {
	MsgId  int64   `json:"msgId"`
	MsgIds []int64 `json:"msgIds"`
	RecId  int64   `json:"recId"`
}
