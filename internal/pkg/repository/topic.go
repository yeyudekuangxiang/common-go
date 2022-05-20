package repository

import (
	"fmt"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	entity2 "mio/internal/pkg/model/entity"
)

var DefaultTopicRepository ITopicRepository = NewTopicRepository(app.DB)

type ITopicRepository interface {
	GetTopicPageList(by GetTopicPageListBy) (list []entity2.Topic, total int64)
	FindById(topicId int64) entity2.Topic
	Save(topic *entity2.Topic) error
	AddTopicLikeCount(topicId int64, num int) error
	GetTopicList(by GetTopicListBy) []entity2.Topic
	GetFlowPageList(by GetTopicFlowPageListBy) (list []entity2.Topic, total int64)
	UpdateColumn(id int64, key string, value interface{}) error
}

func NewTopicRepository(db *gorm.DB) TopicRepository {
	return TopicRepository{DB: db}
}

type TopicRepository struct {
	DB *gorm.DB
}

func (u TopicRepository) GetTopicPageList(by GetTopicPageListBy) (list []entity2.Topic, total int64) {
	list = make([]entity2.Topic, 0)
	//db := u.DB.Model(entity.Topic{})
	db := u.DB.Table("topic").
		Joins("inner join topic_tag on topic.id = topic_tag.topic_id")
	if by.ID > 0 {
		db.Where("topic.id = ?", by.ID)
	}
	if by.TopicTagId > 0 {
		db.Where("topic_tag.tag_id = ?", by.TopicTagId)
	}
	if by.Status > 0 {
		db.Where("status = ?", by.Status)
	}
	db.Select("topic.*").Group("topic.id")

	db2 := u.DB.Table("(?) as t", db)
	err := db2.Count(&total).
		Offset(by.Offset).
		Limit(by.Limit).
		Order("sort desc,updated_at desc,id desc").
		Preload("Tags").
		Find(&list).Error
	if err != nil {
		panic(err)
	}
	return
}
func (u TopicRepository) FindById(topicId int64) entity2.Topic {
	topic := entity2.Topic{}
	err := u.DB.Model(entity2.Topic{}).Where("id = ?", topicId).First(&topic).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return topic
}
func (u TopicRepository) Save(topic *entity2.Topic) error {
	return u.DB.Save(topic).Error
}
func (u TopicRepository) AddTopicLikeCount(topicId int64, num int) error {
	db := u.DB.Model(entity2.Topic{}).
		Where("id = ?", topicId)
	//避免点赞数为负数
	if num < 0 {
		db.Where("like_count >= ?", -num)
	}
	return db.Update("like_count", gorm.Expr("like_count + ?", num)).Error
}
func (u TopicRepository) GetTopicList(by GetTopicListBy) []entity2.Topic {
	list := make([]entity2.Topic, 0)
	db := u.DB.Model(entity2.Topic{})
	if len(by.TopicIds) > 0 {
		db.Where("id in (?)", by.TopicIds)
	}
	if by.Status > 0 {
		db.Where("status = ?", by.Status)
	}
	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return list
}
func (u TopicRepository) GetFlowPageList(by GetTopicFlowPageListBy) (list []entity2.Topic, total int64) {
	list = make([]entity2.Topic, 0)
	db := u.DB.Table(fmt.Sprintf("%s as flow", entity2.TopicFlow{}.TableName())).
		Joins(fmt.Sprintf("inner join %s as topic on flow.topic_id = topic.id", entity2.Topic{}.TableName())).
		Joins(fmt.Sprintf("left join %s as tag on tag.topic_id = topic.id", entity2.TopicTag{}.TableName()))

	if by.UserId > 0 {
		db.Where("flow.user_id = ?", by.UserId)
	}
	if by.TopicTagId > 0 {
		db.Where("tag.tag_id = ?", by.TopicTagId)
	}
	if by.TopicId > 0 {
		db.Where("topic.id = ?", by.TopicId)
	}
	if by.Status > 0 {
		db.Where("topic.status = ?", by.Status)
	}

	db.Select("topic.*,max(flow.sort) as fsort").Group("topic.id")

	db2 := u.DB.Table("(?) as t", db).Order("fsort desc")

	err := db2.Count(&total).
		Offset(by.Offset).
		Limit(by.Limit).
		Order("created_at desc,id desc").
		Preload("Tags").
		Find(&list).Error
	if err != nil {
		panic(err)
	}
	return
}
func (u TopicRepository) UpdateColumn(id int64, key string, value interface{}) error {
	return u.DB.Model(entity2.Topic{}).Where("id = ?", id).UpdateColumn(key, value).Error
}
