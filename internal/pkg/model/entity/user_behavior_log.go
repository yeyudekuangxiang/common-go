package entity

import "time"

type UserBehaviorLog struct {
	Id         int64     `json:"id"`
	Tp         string    `json:"tp"`
	Data       string    `json:"data"`
	Ip         string    `json:"ip"`
	Result     string    `json:"result"`
	ResultCode string    `json:"resultCode"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (UserBehaviorLog) TableName() string {
	return "user_behavior_log"
}
