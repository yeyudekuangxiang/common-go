package message

import (
	"mio/internal/pkg/model/entity"
	"time"
)

const (
	MsgTypeNone    = iota
	MsgTypeLike    //点赞
	MsgTypeComment //评论||回复
	MsgTypeSystem  //系统消息
)

const (
	MsgTurnTypeNone           = iota
	MsgTurnTypeArticle        //跳转文章
	MsgTurnTypeArticleComment //跳转文章评论
	MsgTurnTypeOrder          //跳转订单
	MsgTurnTypeGoods          //跳转商品
	MsgTurnTypeGoodsComment   //跳转商品评论
)

type SendWebMessage struct {
	SendId       int64  `json:"sendId"`
	RecId        int64  `json:"recId"`
	Key          string `json:"key"`
	Type         int    `json:"type"`     // 1点赞 2评论 3回复 4发布 5精选 6违规 7合作社 8 商品 9订单
	TurnType     int    `json:"turnType"` // object_turn_type 1文章 2评论 3订单 4商品
	TurnId       int64  `json:"turnId"`   // object_turn_id   要跳转的object的id
	MessageNotes string `json:"messageNotes"`
}

type SetHaveReadMessage struct {
	MsgIds []string `json:"msgIds" form:"msgIds"`
	RecId  int64    `json:"recId" form:"recId"`
}

type GetWebMessage struct {
	UserId int64    `json:"userId"`
	Status int      `json:"status"`
	Type   int      `json:"type"`
	Types  []string `json:"types"`
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
}

type GetWebMessageResp struct {
	//message
	Id             int64     `json:"id"`
	MessageContent string    `json:"messageContent"`
	MessageNotes   string    `json:"messageNotes"`
	Type           int       `json:"type"`   // 1点赞 2评论 3回复 4发布 5精选 6违规 7合作社
	Status         int       `json:"status"` // 1未读 2已读
	CreatedAt      time.Time `json:"createdAt"`
	//turn obj
	TurnType int    `json:"turnType"` // 1文章 2评论 3订单 4商品
	TurnId   string `json:"turnId"`
	//user
	SendId int64            `json:"sendId"`
	User   entity.ShortUser `json:"user"`
}

type GetWebMessageCount struct {
	RecId int64 `json:"recId"`
}

type GetWebMessageCountResp struct {
	Total               int64 `json:"total"`
	InteractiveMsgTotal int64 `json:"exchangeMsgTotal"`
	SystemMsgTotal      int64 `json:"systemMsgTotal"`
}

type GetMessageTemplate struct {
	Id          int64     `json:"id"`
	Key         string    `json:"key"`
	Type        int       `json:"type"`
	TempContent string    `json:"tempContent"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
