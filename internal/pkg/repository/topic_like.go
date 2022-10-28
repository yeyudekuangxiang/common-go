package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
)

type (
	TopicLikeModel interface {
		Save(like *entity.TopicLike) error
		FindBy(by FindTopicLikeBy) entity.TopicLike
		GetListBy(by GetTopicLikeListBy) []entity.TopicLike
	}

	defaultTopicLikeModel struct {
		ctx *context.MioContext
	}
)

func NewTopicLikeRepository(ctx *context.MioContext) TopicLikeModel {
	return defaultTopicLikeModel{
		ctx: ctx,
	}
}

func (d defaultTopicLikeModel) Save(like *entity.TopicLike) error {
	return d.ctx.DB.Save(like).Error
}
func (d defaultTopicLikeModel) FindBy(by FindTopicLikeBy) entity.TopicLike {
	like := entity.TopicLike{}
	db := d.ctx.DB.Model(like)
	if by.TopicId > 0 {
		db.Where("topic_id = ?", by.TopicId)
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
func (d defaultTopicLikeModel) GetListBy(by GetTopicLikeListBy) []entity.TopicLike {
	list := make([]entity.TopicLike, 0)
	db := d.ctx.DB.Model(entity.TopicLike{})
	if len(by.TopicIds) > 0 {
		db.Where("topic_id in (?)", by.TopicIds)
	}
	if len(by.UserIds) > 0 {
		db.Where("user_id in (?)", by.UserIds)
	}
	if by.TopicId > 0 {
		db.Where("topic_id = ?", by.TopicId)
	}
	if by.UserId > 0 {
		db.Where("user_id = ?", by.UserId)
	}
	if by.Status != 0 {
		db.Where("status = ?", by.Status)
	}
	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return list
}
