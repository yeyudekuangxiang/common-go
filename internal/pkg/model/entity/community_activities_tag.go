package entity

import (
	"time"
)

// CommunityActivitiesTag 社区用户举办的活动表标签
type CommunityActivitiesTag struct {
	Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (CommunityActivitiesTag) TableName() string {
	return "community_activities_tag"
}
