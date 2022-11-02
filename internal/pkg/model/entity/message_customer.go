package entity

import "time"

type MessageCustomer struct {
	Id        int64     `json:"id"`
	RecId     int64     `json:"recId"`
	MessageId int64     `json:"messageId"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (m MessageCustomer) TableName() string {
	return "message_customer"
}

type UserWebMessage struct {
	Id             int64     `json:"id"`
	RecId          int64     `json:"recId"`
	MessageId      int64     `json:"messageId"`
	MessageContent string    `json:"messageContent"`
	Status         int       `json:"status"` //1未读 2已读
	Type           int       `json:"type"`   // 1互动消息 2系统消息
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type UserWebMessageV2 struct {
	MessageId      int64     `json:"messageId"`
	MessageContent string    `json:"messageContent"`
	CreatedAt      time.Time `json:"createdAt"`
}
