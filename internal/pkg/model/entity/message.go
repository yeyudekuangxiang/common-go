package entity

import "time"

type Message struct {
	Id        int64     `json:"id"`
	SendId    int64     `json:"sendId"`
	RecId     int64     `json:"recId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (m Message) TableName() string {
	return "message"
}
