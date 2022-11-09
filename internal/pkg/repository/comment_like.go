package repository

import (
	"gorm.io/gorm"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
)

type (
	CommentLikeModel interface {
		Save(like *entity.CommentLike) error
		Create(commentLike *entity.CommentLike) error
		Update(commentLike *entity.CommentLike) error
		FindBy(by FindCommentLikeBy) entity.CommentLike
		GetListBy(by GetCommentLikeListBy) []entity.CommentLike
	}

	defaultCommentLikeModel struct {
		ctx *mioContext.MioContext
	}
)

func NewCommentLikeRepository(ctx *mioContext.MioContext) CommentLikeModel {
	return defaultCommentLikeModel{
		ctx: ctx,
	}
}

type CommentLikeRepository struct {
	ctx *mioContext.MioContext
}

func (m defaultCommentLikeModel) Save(like *entity.CommentLike) error {
	return m.ctx.DB.Save(like).Error
}

func (m defaultCommentLikeModel) Create(commentLike *entity.CommentLike) error {
	return m.ctx.DB.Debug().Create(commentLike).Error
}

func (m defaultCommentLikeModel) Update(commentLike *entity.CommentLike) error {
	return m.ctx.DB.Debug().Updates(commentLike).Error
}

func (m defaultCommentLikeModel) FindBy(by FindCommentLikeBy) entity.CommentLike {
	like := entity.CommentLike{}
	db := m.ctx.DB.Model(&like)
	if by.CommentId > 0 {
		db.Where("comment_id = ?", by.CommentId)
	}
	if by.UserId > 0 {
		db.Where("user_id = ?", by.UserId)
	}
	if err := db.First(&like).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return like
}
func (m defaultCommentLikeModel) GetListBy(by GetCommentLikeListBy) []entity.CommentLike {
	list := make([]entity.CommentLike, 0)
	db := m.ctx.DB.Model(&entity.CommentLike{})
	if len(by.CommentIds) > 0 {
		db.Where("comment_id in (?)", by.CommentIds)
	}
	if by.UserId > 0 {
		db.Where("user_id = ?", by.UserId)
	}
	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return list
}
