package entity

import (
	"mio/internal/pkg/model"
)

type TopicFlow struct {
	ID             int64
	UserId         int64
	TopicId        int64
	SeeCount       int
	ShowCount      int
	Sort           int
	TopicCreatedAt model.Time
	TopicUpdatedAt model.Time
}

func (TopicFlow) TableName() string {
	return "topic_flow"
}
