package kumiaoCommunity

import (
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
