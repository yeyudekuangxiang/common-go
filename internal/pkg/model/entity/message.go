package entity

import "time"

type Message struct {
	Id        int64     `json:"id"`
	SendId    int64     `json:"sendId"`
	RecId     int64     `json:"recId"`
	Type      int       `json:"type"` // 1 互动消息 2 系统消息
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (m Message) TableName() string {
	return "message"
}
