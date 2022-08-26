package repository

import (
	"fmt"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultTopicRepository ITopicRepository = NewTopicRepository(app.DB)

type ITopicRepository interface {
	GetTopicPageList(by GetTopicPageListBy) (list []entity.Topic, total int64)
	FindById(topicId int64) entity.Topic
	Save(topic *entity.Topic) error
	//FindTopicOneBy() entity.Topic
	AddTopicLikeCount(topicId int64, num int) error
	GetTopicList(by FindTopicBy) []entity.Topic
	GetFlowPageList(by GetTopicFlowPageListBy) (list []entity.Topic, total int64)
	UpdateColumn(id int64, key string, value interface{}) error
}

func NewTopicRepository(db *gorm.DB) TopicRepository {
	return TopicRepository{DB: db}
}

type TopicRepository struct {
	DB *gorm.DB
}

func (u TopicRepository) GetTopicPageList(by GetTopicPageListBy) (list []entity.Topic, total int64) {
	list = make([]entity.Topic, 0)

	db := u.DB.Table("topic")

	if by.ID > 0 {
		db.Where("topic.id = ?", by.ID)
	}
	if by.TopicTagId > 0 {
		db.Joins("inner join topic_tag on topic.id = topic_tag.topic_id").Where("topic_tag.tag_id = ?", by.TopicTagId)
	}
	if by.Status > 0 {
		db.Where("topic.status = ?", by.Status)
	}
	err := db.Count(&total).
		Offset(by.Offset).
		Limit(by.Limit).
		Order("is_top desc, is_essence desc, sort desc,updated_at desc,id desc").
		Preload("Tags").
		Find(&list).Error
	if err != nil {
		panic(err)
	}
	return
}

func (u TopicRepository) FindById(topicId int64) entity.Topic {
	topic := entity.Topic{}
	err := u.DB.Model(entity.Topic{}).Preload("User").Where("id = ?", topicId).First(&topic).Error
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

func (u TopicRepository) GetTopicList(by FindTopicBy) []entity.Topic {
	list := make([]entity.Topic, 0)
	db := u.DB.Model(entity.Topic{})
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
func (u TopicRepository) GetFlowPageList(by GetTopicFlowPageListBy) (list []entity.Topic, total int64) {
	list = make([]entity.Topic, 0)
	db := u.DB.Table(fmt.Sprintf("%s as flow", entity.TopicFlow{}.TableName())).
		Joins(fmt.Sprintf("inner join %s as topic on flow.topic_id = topic.id", entity.Topic{}.TableName())).
		Where("flow.user_id = ?", by.UserId)

	if by.TopicTagId > 0 {
		db.Joins(fmt.Sprintf("left join %s as tag on tag.topic_id = topic.id", entity.TopicTag{}.TableName())).
			Where("tag.tag_id = ?", by.TopicTagId)
	}
	if by.TopicId > 0 {
		db.Where("topic.id = ?", by.TopicId)
	}
	if by.Status > 0 {
		db.Where("topic.status = ?", by.Status)
	}

	err := db.Select("topic.*,flow.sort fsort").
		Count(&total).
		Offset(by.Offset).
		Limit(by.Limit).
		Order("fsort desc,flow.created_at desc,flow.id desc").
		Preload("Tags").
		Find(&list).Error
	if err != nil {
		panic(err)
	}
	return
}
func (u TopicRepository) UpdateColumn(id int64, key string, value interface{}) error {
	return u.DB.Model(entity.Topic{}).Where("id = ?", id).UpdateColumn(key, value).Error
}
