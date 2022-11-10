package message

import (
	"mio/internal/pkg/model/entity"
	"time"
)

type SendWebMessage struct {
	SendId   int64  `json:"sendId"`
	RecId    int64  `json:"recId"`
	Key      string `json:"key"`
	Type     int    `json:"type"`     // 1点赞 2评论 3回复 4发布 5精选 6违规 7合作社 8 商品 9订单
	TurnType int    `json:"turnType"` // object_turn_type 1文章 2评论 3订单 4商品
	TurnId   int64  `json:"turnId"`   // object_turn_id   要跳转的object的id
	ShowId   int64  `json:"forId"`    // object_show_id   展示id
	//ShowParentId int64  `json:"showParentId"`
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
	Type           int       `json:"type"`   // 1点赞 2评论 3回复 4发布 5精选 6违规 7合作社
	Status         int       `json:"status"` // 1未读 2已读
	CreatedAt      time.Time `json:"createdAt"`
	//turn obj
	TurnType  int    `json:"turnType"` // 1文章 2评论 3订单 4商品
	TurnId    int64  `json:"turnId"`
	ShowId    int64  `json:"showId"`
	TurnNotes string `json:"turnNotes"` // 简介
	//user
	SendId int64            `json:"sendId"`
	User   entity.ShortUser `json:"user"`
}

type GetWebMessageCount struct {
	RecId int64 `json:"recId"`
}

type GetWebMessageCountResp struct {
	Total            int64 `json:"total"`
	ExchangeMsgTotal int64 `json:"exchangeMsgTotal"`
	SystemMsgTotal   int64 `json:"systemMsgTotal"`
}
