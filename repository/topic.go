package repository

import (
	"mio/core/app"
	"mio/model/entity"
)

var DefaultTopicRepository ITopicRepository = NewTopicRepository()

type ITopicRepository interface {
	GetTopicPageList(by GetTopicPageListBy) (list []entity.Topic, total int64)
}

func NewTopicRepository() TopicRepository {
	return TopicRepository{}
}

type TopicRepository struct {
}

func (u TopicRepository) GetTopicPageList(by GetTopicPageListBy) (list []entity.Topic, total int64) {
	list = make([]entity.Topic, 0)
	//db := app.DB.Model(entity.Topic{})
	db := app.DB.Table("topic").
		Joins("inner join topic_tag on topic.id = topic_tag.topic_id")
	if by.ID > 0 {
		db.Where("topic.id = ?", by.ID)
	}
	if by.TopicTagId > 0 {
		db.Where("topic_tag.tag_id = ?", by.TopicTagId)
	}
	db.Select("topic.*").Group("topic.id")

	db2 := app.DB.Table("(?) as t", db)
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
