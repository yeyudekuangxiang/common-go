package kumiaoCommunity

import (
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"time"
)

type TopicChangeLikeResp struct {
	TopicTitle  string `json:"topic"`
	TopicId     int64  `json:"topicId"`
	TopicUserId int64  `json:"TopicUserId"`
	LikeStatus  int    `json:"likeStatus"`
	IsFirst     bool   `json:"isFirst"`
}

type CommentChangeLikeResp struct {
	CommentMessage string `json:"topic"`
	CommentId      int64  `json:"topicId"`
	CommentUserId  int64  `json:"TopicUserId"`
	LikeStatus     int    `json:"likeStatus"`
	IsFirst        bool   `json:"isFirst"`
}

type ChangeUserPosition struct {
	UserId       int64  `json:"userId"`
	Position     string `json:"position"`
	PositionIcon string `json:"positionIcon"`
}

type ChangeUserPartner struct {
	UserId  int64 `json:"userId"`
	Partner int   `json:"partner"`
}

type CommentCount struct {
	Date    time.Time
	TopicId int64
	Total   int64
}

type ChangeUserState struct {
	UserId int64 `json:"userId"`
	Status int   `json:"state"`
}

type TrackOrderZhuGe struct {
	OpenId        string
	CertificateId string
	ProductItemId string
	OrderId       string
	Partnership   entity.PartnershipType
	Title         string
	CateTitle     string
}

type TrackLoginZhuGe struct {
	CateTitle string
}

type TopicDetail struct {
	entity.Topic
	IsLike        bool             `json:"isLike"`
	UpdatedAtDate string           `json:"updatedAtDate"` //03-01
	User          entity.ShortUser `json:"user"`
}

type TurnCommentReq struct {
	UserId int64  `json:"userId"`
	Types  int    `json:"types"`
	TurnId string `json:"turnId"`
}

type APIComment struct {
	Id            string           `gorm:"primaryKey;autoIncrement" json:"id"`
	ObjId         string           `gorm:"type:int8;not null" json:"objId" json:"objId"` // 对象id （文章）
	ObjType       int8             `gorm:"type:int2" json:"objType"`                     // 保留字段
	Message       string           `gorm:"type:text;not null" json:"message"`
	MemberId      int64            `gorm:"type:int8;not null" json:"memberId"`                // 评论用户id
	RootCommentId int64            `gorm:"type:int8;not null:default:0" json:"rootCommentId"` // 根评论id，不为0表示是回复评论
	ToCommentId   int64            `gorm:"type:int8;not null:default:0" json:"toCommentId"`   // 父评论id，为0表示是根评论
	Floor         int32            `gorm:"type:int4;not null:default:0" json:"floor"`         // 评论楼层
	Count         int32            `gorm:"type:int4;not null:default:0" json:"count"`         // 该评论下评论总数
	RootCount     int32            `gorm:"type:int4;not null:default:0" json:"rootCount"`     // 该评论下根评论总数
	LikeCount     int32            `gorm:"type:int4;not null:default:0" json:"likeCount"`     // 该评论点赞总数
	HateCount     int32            `gorm:"type:int4;not null:default:0" json:"hateCount"`     // 该评论点踩总数
	State         int8             `gorm:"type:int2;not null:default:0" json:"state"`         // 状态 0-正常 1-隐藏
	DelReason     string           `gorm:"type:varchar(200)" json:"delReason"`                //删除理由
	Attrs         int8             `gorm:"type:int2" json:"attrs"`                            // 属性 00-正常 10-运营置顶 01-用户置顶 保留字段
	Version       int64            `gorm:"type:int8;version" json:"version"`                  // 版本号 保留字段
	CreatedAt     model.Time       `gorm:"createdAt" json:"createdAt"`
	UpdatedAt     model.Time       `gorm:"updatedAt" json:"updatedAt"`
	Member        entity.ShortUser `gorm:"foreignKey:ID;references:MemberId" json:"member"` // 评论用户
	IsAuthor      int8             `gorm:"type:int2" json:"isAuthor"`                       // 是否作者
	RootChild     []*APIComment    `gorm:"foreignKey:RootCommentId;association_foreignKey:Id" json:"rootChild"`
}

type APICommentResp struct {
	Id        string            `gorm:"primaryKey;autoIncrement" json:"id"`
	Message   string            `gorm:"type:text;not null" json:"message"`
	MemberId  int64             `gorm:"type:int8;not null" json:"memberId"`            // 评论用户id
	Floor     int32             `gorm:"type:int4;not null:default:0" json:"floor"`     // 评论楼层
	Count     int32             `gorm:"type:int4;not null:default:0" json:"count"`     // 该评论下评论总数
	RootCount int32             `gorm:"type:int4;not null:default:0" json:"rootCount"` // 该评论下根评论总数
	LikeCount int32             `gorm:"type:int4;not null:default:0" json:"likeCount"` // 该评论点赞总数
	HateCount int32             `gorm:"type:int4;not null:default:0" json:"hateCount"` // 该评论点踩总数
	CreatedAt model.Time        `json:"createdAt"`
	UpdatedAt model.Time        `json:"updatedAt"`
	Member    entity.ShortUser  `json:"member"`
	Detail    Detail            `json:"detail"`
	IsAuthor  int8              `json:"isAuthor"` // 是否作者
	IsLike    int               `json:"isLike"`
	RootChild []*APICommentResp `json:"rootChild"`
}

type Detail struct {
	ObjId       string `json:"objId"`
	ObjType     int64  `json:"objType"`
	ImageList   string `json:"imageList"`
	Description string `json:"description"`
}

func (a APIComment) ApiComment() *APICommentResp {
	return &APICommentResp{
		Id:        a.Id,
		Message:   a.Message,
		MemberId:  a.MemberId,
		Floor:     a.Floor,
		Count:     a.Count,
		RootCount: a.RootCount,
		LikeCount: a.LikeCount,
		HateCount: a.HateCount,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		Member:    a.Member,
		IsAuthor:  a.IsAuthor,
	}
}
