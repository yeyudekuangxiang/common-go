package repository

import (
	"gorm.io/gorm"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
)

type (
	CarbonCommentLikeModel interface {
		Save(like *entity.CarbonCommentLike) error
		Create(commentLike *entity.CarbonCommentLike) error
		Update(commentLike *entity.CarbonCommentLike) error
		FindBy(by FindCommentLikeBy) entity.CarbonCommentLike
		GetListBy(by GetCommentLikeListBy) []entity.CarbonCommentLike
	}

	defaultCarbonCommentLikeModel struct {
		ctx *mioContext.MioContext
	}
)

func NewCarbonCommentLikeRepository(ctx *mioContext.MioContext) CarbonCommentLikeModel {
	return defaultCarbonCommentLikeModel{
		ctx: ctx,
	}
}

type CarbonCommentLikeRepository struct {
	ctx *mioContext.MioContext
}

func (m defaultCarbonCommentLikeModel) Save(like *entity.CarbonCommentLike) error {
	return m.ctx.DB.Save(like).Error
}

func (m defaultCarbonCommentLikeModel) Create(commentLike *entity.CarbonCommentLike) error {
	return m.ctx.DB.Debug().Create(commentLike).Error
}

func (m defaultCarbonCommentLikeModel) Update(commentLike *entity.CarbonCommentLike) error {
	return m.ctx.DB.Debug().Updates(commentLike).Error
}

func (m defaultCarbonCommentLikeModel) FindBy(by FindCommentLikeBy) entity.CarbonCommentLike {
	like := entity.CarbonCommentLike{}
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
func (m defaultCarbonCommentLikeModel) GetListBy(by GetCommentLikeListBy) []entity.CarbonCommentLike {
	list := make([]entity.CarbonCommentLike, 0)
	db := m.ctx.DB.Model(&entity.CarbonCommentLike{})
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
