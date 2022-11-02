package message

import (
	"mio/internal/pkg/model/entity"
	"time"
)

type SendWebMessage struct {
	SendId   int64  `json:"sendId"`
	RecId    int64  `json:"recId"`
	Key      string `json:"key"`
	RecObjId int64  `json:"recObjId"`
	Type     int    `json:"type" default:"1"`
}

type SetHaveReadMessage struct {
	MsgId  int64   `json:"msgId" form:"msgId"`
	MsgIds []int64 `json:"msgIds" form:"msgIds"`
	RecId  int64   `json:"recId" form:"recId"`
}

type GetWebMessage struct {
	UserId int64 `json:"userId"`
	Status int   `json:"status"`
	Type   int   `json:"type"`
	Types  []int `json:"types"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

type GetWebMessageResp struct {
	//message
	Id             int64     `json:"id"`
	MessageContent string    `json:"messageContent"`
	Type           int       `json:"type"` //1点赞 2评论 3回复 4发布 5精选 6违规 7合作社
	CreatedAt      time.Time `json:"createdAt"`
	//turn obj
	TurnType  int    `json:"TurnType"` // 1文章 2评论 3订单 4商品
	TurnId    int64  `json:"TurnId"`
	TurnNotes string `json:"turnNotes"` // 简介
	//user
	User entity.ShortUser
}

type GetWebMessageCount struct {
	RecId int64 `json:"recId"`
}

type GetWebMessageCountResp struct {
	Total            int64 `json:"total"`
	ExchangeMsgTotal int64 `json:"exchangeMsgTotal"`
	SystemMsgTotal   int64 `json:"systemMsgTotal"`
}
