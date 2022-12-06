package message

import "time"

type SendMessage struct {
	SendId       int64  `json:"sendId"`
	RecId        int64  `json:"recId"`
	Type         int    `json:"type"`     // 1点赞 2评论 3回复 4发布 5精选 6违规 7合作社 8 商品 9订单
	TurnType     int    `json:"turnType"` // 1文章 2评论 3订单 4商品
	TurnId       int64  `json:"turnId"`
	Message      string `json:"message"`
	MessageNotes string `json:"messageNotes"`
}

type FindMessageParams struct {
	MessageIds []int64   `json:"messageIds"`
	MessageId  int64     `json:"messageId"`
	SendId     int64     `json:"sendId"`
	RecId      int64     `json:"recId"`
	Type       int       `json:"type"` // 1互动消息 2酷喵圈社区 3二手交易
	Types      []string  `json:"types"`
	Status     int       `json:"status" default:"1"` // 1未读 2已读
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	Limit      int       `json:"limit"`
	Offset     int       `json:"offset"`
}

type SetHaveReadMessageParams struct {
	MsgIds []string `json:"msgIds"`
	RecId  int64    `json:"recId"`
}

type IMSendMsg struct {
}

type IMGetMsg struct {
}
