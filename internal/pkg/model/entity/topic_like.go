package entity

import (
	"mio/internal/pkg/model"
)

type TopicLike struct {
	Id        int        `json:"id"`
	TopicId   int        `json:"topicId"`
	UserId    int        `json:"userId"`
	Status    int8       `json:"status"`
	CreatedAt model.Time `json:"createdAt"`
	UpdatedAt model.Time `json:"updatedAt"`
}

func (TopicLike) TableName() string {
	return "topic_like"
}
