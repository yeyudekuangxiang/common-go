package entity

import (
	"time"
)

type BdSceneUser struct {
	ID             int64     `json:"ID"`
	PlatformKey    string    `json:"platformKey"`
	PlatformUserId string    `json:"platformUserId,omitempty"` //外站用户id
	Phone          string    `json:"phone,omitempty"`          //外站用户手机
	OpenId         string    `json:"openId,omitempty"`         //本站用户openId
	UnionId        string    `json:"unionId,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func (BdSceneUser) TableName() string {
	return "bd_scene_user"
}
