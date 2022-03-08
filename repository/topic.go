package repository

import (
	"gorm.io/gorm"
	"mio/core/app"
	"mio/model/entity"
)

var DefaultTopicRepository ITopicRepository = NewTopicRepository(app.DB)

type ITopicRepository interface {
	GetTopicPageList(by GetTopicPageListBy) (list []entity.Topic, total int64)
	FindById(topicId int64) entity.Topic
	Save(topic *entity.Topic) error
	AddTopicLikeCount(topicId int64, num int) error
}

func NewTopicRepository(db *gorm.DB) TopicRepository {
	return TopicRepository{DB: db}
}

type TopicRepository struct {
	DB *gorm.DB
}

func (u TopicRepository) GetTopicPageList(by GetTopicPageListBy) (list []entity.Topic, total int64) {
	list = make([]entity.Topic, 0)
	//db := u.DB.Model(entity.Topic{})
	db := u.DB.Table("topic").
		Joins("inner join topic_tag on topic.id = topic_tag.topic_id")
	if by.ID > 0 {
		db.Where("topic.id = ?", by.ID)
	}
	if by.TopicTagId > 0 {
		db.Where("topic_tag.tag_id = ?", by.TopicTagId)
	}
	db.Select("topic.*").Group("topic.id")

	db2 := u.DB.Table("(?) as t", db)
	err := db2.Count(&total).
		Offset(by.Offset).
		Limit(by.Limit).
		Order("sort desc").
		Preload("Tags").
		Find(&list).Error
	if err != nil {
		panic(err)
	}
	return
}
func (u TopicRepository) FindById(topicId int64) entity.Topic {
	topic := entity.Topic{}
	err := u.DB.Model(entity.Topic{}).Where("id = ?", topicId).First(&topic).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return topic
}
func (u TopicRepository) Save(topic *entity.Topic) error {
	return u.DB.Save(topic).Error
}
func (u TopicRepository) AddTopicLikeCount(topicId int64, num int) error {
	db := u.DB.Model(entity.Topic{}).
		Where("id = ?", topicId)
	//避免点赞数为负数
	if num < 0 {
		db.Where("like_count >= ?", -num)
	}
	return db.Update("like_count", gorm.Expr("like_count + ?", num)).Error
}
