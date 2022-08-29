package entity

import (
	"mio/internal/pkg/model"
)

type CommentLike struct {
	Id        int64      `json:"id"`
	CommentId int64      `json:"commentId"`
	UserId    int64      `json:"userId"`
	Status    int8       `json:"status"`
	CreatedAt model.Time `json:"createdAt"`
	UpdatedAt model.Time `json:"updatedAt"`
}

func (CommentLike) TableName() string {
	return "comment_like"
}
