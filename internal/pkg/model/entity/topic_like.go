package entity

import (
	"mio/internal/pkg/model"
)

type TopicLike struct {
	Id        int64      `json:"id"`
	TopicId   int64      `json:"topicId"`
	UserId    int64      `json:"userId"`
	Status    int8       `json:"status"`
	CreatedAt model.Time `json:"createdAt"`
	UpdatedAt model.Time `json:"updatedAt"`
}

func (TopicLike) TableName() string {
	return "topic_like"
}
