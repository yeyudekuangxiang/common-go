package entity

import "time"

type MessageContent struct {
	MessageId      int64     `json:"messageId"`
	MessageContent string    `json:"messageContent"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func (m MessageContent) TableName() string {
	return "message_content"
}
