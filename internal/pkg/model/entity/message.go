package entity

import "time"

type Message struct {
	Id       int64 `json:"id"`
	SendId   int64 `json:"sendId"`
	RecId    int64 `json:"recId"`
	Type     int   `json:"type"`     // 1点赞 2评论 3回复 4发布 5精选 6违规 7合作社
	TurnType int   `json:"turnType"` // 1文章 2评论 3订单 4商品
	TurnId   int64 `json:"turnId"`
	//ShowId    int64     `json:"showId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (m Message) TableName() string {
	return "message"
}

type UserWebMessage struct {
	//message
	Id             int64     `json:"id"`
	MessageContent string    `json:"messageContent"`
	Type           int       `json:"type"`   //1点赞 2评论 3回复 4发布 5精选 6违规 7合作社
	Status         int       `json:"status"` //1未读 2已读
	CreatedAt      time.Time `json:"createdAt"`
	//obj
	TurnType int   `json:"turnType"` // 1文章 2 评论 3订单 4商品
	TurnId   int64 `json:"turnId"`
	ShowId   int64 `json:"showId"`
	//user
	SendId int64 `json:"sendId"`
	//NickName  string `json:"nickName"`
	//AvatarUrl string `json:"avatarUrl"`
}
