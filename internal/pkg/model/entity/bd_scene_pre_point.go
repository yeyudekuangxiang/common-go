package entity

import (
	"time"
)

type BdScenePrePoint struct {
	ID             int64     `json:"ID"`
	PlatformKey    string    `json:"platformKey"`
	PlatformUserId string    `json:"platformUserId,omitempty"` //外站用户id
	Point          string    `json:"phone,omitempty"`          //外站用户手机
	OpenId         string    `json:"openId,omitempty"`         //本站用户openId
	Status         int       `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func (BdScenePrePoint) TableName() string {
	return "bd_scene_pre_point"
}
