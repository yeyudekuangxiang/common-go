package repository

import (
	"mio/core/app"
	"mio/model"
)

var DefaultTopicRepository ITopicRepository = NewTopicRepository()

type ITopicRepository interface {
	GetTopicPageList(by GetTopicPageListBy) (list []model.Topic, total int64)
}

func NewTopicRepository() TopicRepository {
	return TopicRepository{}
}

type TopicRepository struct {
}

func (u TopicRepository) GetTopicPageList(by GetTopicPageListBy) (list []model.Topic, total int64) {
	list = make([]model.Topic, 0)
	db := app.DB.Model(model.Topic{})
	if by.ID > 0 {
		db.Where("id = ?", by.ID)
	}
	if by.TopicTagId > 0 {
		db.Where("topic_tag_id = ?", by.TopicTagId)
	}
	err := db.Count(&total).Offset(by.Offset).Limit(by.Limit).Order("sort desc").Find(&list).Error
	if err != nil {
		panic(err)
	}
	return
}
