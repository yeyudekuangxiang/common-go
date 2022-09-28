package entity

import "time"

type BdSceneCallback struct {
	ID             int
	PlatformKey    string
	PlatformUserId string
	OpenId         string
	BizId          string
	SourceKey      string
	Body           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (BdSceneCallback) TableName() string {
	return "bd_scene_callback"
}
