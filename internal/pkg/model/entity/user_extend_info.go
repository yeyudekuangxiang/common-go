package entity

import "time"

type UserExtendInfo struct {
	ID        int
	Openid    string    `json:"openid"`
	Uid       int64     `json:"uid"`
	Ip        string    `json:"ip"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (UserExtendInfo) TableName() string {
	return "user_extend_info"
}
