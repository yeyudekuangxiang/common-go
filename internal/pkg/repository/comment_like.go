package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultCommentLikeRepository = NewCommentLikeRepository(app.DB)

func NewCommentLikeRepository(db *gorm.DB) CommentLikeRepository {
	return CommentLikeRepository{
		DB: db,
	}
}

type CommentLikeRepository struct {
	DB *gorm.DB
}

func (t CommentLikeRepository) Save(like *entity.CommentLike) error {
	return t.DB.Save(like).Error
}
func (t CommentLikeRepository) FindBy(by FindCommentLikeBy) entity.CommentLike {
	like := entity.CommentLike{}
	db := t.DB.Model(like)
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
func (t CommentLikeRepository) GetListBy(by GetCommentLikeListBy) []entity.CommentLike {
	list := make([]entity.CommentLike, 0)
	db := t.DB.Model(entity.CommentLike{})
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
