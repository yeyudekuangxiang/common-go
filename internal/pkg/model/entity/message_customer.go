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
