package entity

import "mio/internal/pkg/model"

type EventParticipation struct {
	ID        int64  `json:"id"`
	EventId   string `json:"event_id"`
	Nickname  string `json:"nickname" gorm:"column:nick_name"`
	AvatarUrl string `json:"avatarUrl" gorm:"avatarUrl"`
	Count     int    `json:"count"`
	Time      model.Time
}
