package message

import "time"

type IM struct {
	Id         int64     `json:"id"`
	ChannelId  string    `json:"channelId"`
	SendUserId int64     `json:"sendUserId"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (m IM) TableName() string {
	return "im_message"
}
