package event

import (
	"mio/internal/pkg/model/entity"
	"time"
)

const OrderByEventParticipationCountDesc entity.OrderBy = "order_by_event_participation_count_desc"
const OrderByEventParticipationTimeDesc entity.OrderBy = "order_by_event_participation_time_desc"

type EventParticipation struct {
	ID        int64  `json:"id"`
	EventId   string `json:"event_id"`
	Nickname  string `json:"nickname" gorm:"column:nick_name"`
	AvatarUrl string `json:"avatarUrl" gorm:"avatarUrl"`
	Count     int    `json:"count"`
	Time      time.Time
}
