package entity

import "mio/model"

type TopicUserFlow struct {
	ID             int64
	UserId         int64
	TopicId        int64
	SeeCount       int
	ShowCount      int
	Sort           int
	TopicCreatedAt model.Time
	TopicUpdatedAt model.Time
}

func (TopicUserFlow) TableName() string {
	return "topic_user_flow"
}
