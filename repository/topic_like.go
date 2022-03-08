package repository

import (
	"gorm.io/gorm"
	"mio/model/entity"
)

type TopicLikeRepository struct {
	DB *gorm.DB
}

func (t TopicLikeRepository) Save(like *entity.TopicLike) error {
	return t.DB.Save(like).Error
}
func (t TopicLikeRepository) FindBy(by FindTopicLikeBy) entity.TopicLike {
	like := entity.TopicLike{}
	db := t.DB.Model(like)
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
func (t TopicLikeRepository) GetListBy(by GetTopicLikeListBy) []entity.TopicLike {
	list := make([]entity.TopicLike, 0)
	db := t.DB.Model(entity.TopicLike{})
	if len(by.TopicIds) > 0 {
		db.Where("topic_id in ?", by.TopicIds)
	}
	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return list
}
