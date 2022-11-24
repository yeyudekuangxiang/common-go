package message

import "time"

type MessageTemplate struct {
	Id          int64     `json:"id"`
	Key         string    `json:"key"`
	Type        int       `json:"type"`
	TempContent string    `json:"tempContent"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (m MessageTemplate) TableName() string {
	return "message_template"
}
