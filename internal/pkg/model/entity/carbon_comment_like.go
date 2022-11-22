package entity

import (
	"mio/internal/pkg/model"
)

type CarbonCommentLike struct {
	Id        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	CommentId int64      `gorm:"type:int8;not null" json:"commentId"`
	UserId    int64      `gorm:"type:int8;not null" json:"userId"`
	Status    int8       `gorm:"type:int2;not null" json:"status"`
	CreatedAt model.Time `json:"createdAt"`
	UpdatedAt model.Time `json:"updatedAt"`
}

func (CarbonCommentLike) TableName() string {
	return "carbon_comment_like"
}
