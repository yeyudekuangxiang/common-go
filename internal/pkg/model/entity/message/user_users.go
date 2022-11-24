package message

import "time"

type UserUsers struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"userId"`
	HasUserId int64     `json:"hasUserId"`
	ChannelId string    `json:"channelId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (m UserUsers) TableName() string {
	return "user_users"
}
