package entity

import "mio/internal/pkg/model"

// CommunityActivitiesTag 社区用户举办的活动表标签
type CommunityActivitiesTag struct {
	Id        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string     `json:"title"`
	CreatedAt model.Time `json:"createdAt"`
	UpdatedAt model.Time `json:"updatedAt"`
}

func (CommunityActivitiesTag) TableName() string {
	return "community_activities_tag"
}
