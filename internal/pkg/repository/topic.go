package repository

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
)

type (
	TopicModel interface {
		GetTopicPageList(by GetTopicPageListBy) (list []entity.Topic, total int64)
		FindById(topicId int64) entity.Topic
		Save(topic *entity.Topic) error
		AddTopicLikeCount(topicId int64, num int) error
		GetFlowPageList(by GetTopicFlowPageListBy) (list []entity.Topic, total int64)
		UpdateColumn(id int64, key string, value interface{}) error
		GetMyTopic(by GetTopicPageListBy) ([]*entity.Topic, int64, error)
		GetTopicList(by GetTopicPageListBy) ([]*entity.Topic, int64, error)
		ChangeTopicCollectionCount(id int64, column string, incr int) error
		Trans(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error
		GetTopicNotes(topicIds []int64) []*entity.Topic
	}

	defaultTopicRepository struct {
		ctx *mioContext.MioContext
	}
)

func (d defaultTopicRepository) GetTopicNotes(topicIds []int64) []*entity.Topic {
	topList := make([]*entity.Topic, 0)
	err := d.ctx.DB.Model(&entity.Topic{}).
		Preload("User").
		Preload("Tags").
		Where("topic.id in ?", topicIds).
		Where("topic.status = ?", entity.TopicStatusPublished).
		Group("topic.id").
		Find(&topList).Error
	if err != nil {
		app.Logger.Error(err)
		return []*entity.Topic{}
	}
	return topList
}

func (d defaultTopicRepository) Trans(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return d.ctx.DB.Transaction(fc, opts...)
}

func (d defaultTopicRepository) ChangeTopicCollectionCount(id int64, column string, incr int) error {
	return d.ctx.DB.Model(&entity.Topic{}).Where("id = ?", id).Update(column, gorm.Expr(column+"+?", incr)).Error
}

func (d defaultTopicRepository) GetMyTopic(by GetTopicPageListBy) ([]*entity.Topic, int64, error) {
	topList := make([]*entity.Topic, 0)
	var total int64
	query := d.ctx.DB.Model(&entity.Topic{}).
		Preload("User").
		Preload("Tags").
		Preload("Comment", func(db *gorm.DB) *gorm.DB {
			return db.Where("comment_index.to_comment_id = ?", 0).
				Order("like_count desc").Limit(10)
		}).
		Preload("Comment.RootChild", func(db *gorm.DB) *gorm.DB {
			return db.Where("(select count(*) from comment_index index where index.root_comment_id = comment_index.root_comment_id and index.id <= comment_index.id) <= ?", 3).
				Order("comment_index.like_count desc")
		}).
		Preload("Comment.RootChild.Member").
		Preload("Comment.Member")
	if by.Status != 0 {
		query.Where("topic.status = ?", by.Status)
	}
	err := query.Where("topic.user_id = ?", by.UserId).
		Count(&total).
		Group("topic.id").
		Order("id desc").
		Limit(by.Limit).
		Offset(by.Offset).
		Find(&topList).Error

	if err != nil {
		return nil, 0, err
	}

	return topList, total, nil
}

func (d defaultTopicRepository) GetTopicList(by GetTopicPageListBy) ([]*entity.Topic, int64, error) {
	topList := make([]*entity.Topic, 0)
	var total int64
	query := d.ctx.DB.Model(&entity.Topic{}).
		Preload("User").
		Preload("Tags").
		Preload("Comment", func(db *gorm.DB) *gorm.DB {
			return db.Where("comment_index.to_comment_id = ?", 0).
				Order("like_count desc").Limit(10)
		}).
		Preload("Comment.RootChild", func(db *gorm.DB) *gorm.DB {
			return db.Where("(select count(*) from comment_index index where index.root_comment_id = comment_index.root_comment_id and index.id <= comment_index.id) <= ?", 3).
				Order("comment_index.like_count desc")
		}).
		Preload("Comment.RootChild.Member").
		Preload("Comment.Member")

	if by.ID != 0 {
		query.Where("topic.id = ?", by.ID)
	} else if len(by.Ids) > 0 {
		query.Where("topic.id in ?", by.Ids)
	}

	if by.TopicTagId != 0 {
		query.Joins("inner join topic_tag on topic.id = topic_tag.topic_id").Where("topic_tag.tag_id = ?", by.TopicTagId)
	}
	if by.UserId != 0 {
		query.Where("topic.user_id = ?", by.UserId)
	}
	if by.Status != 0 {
		query.Where("topic.status = ?", by.Status)
	} else {
		query.Where("topic.status = ?", entity.TopicStatusPublished)
	}
	query = query.Count(&total).
		Group("topic.id")
	if by.Order == "time" {
		query.Order("topic.created_at desc, topic.like_count desc, topic.see_count desc, topic.id desc")
	} else if by.Order == "recommend" {
		query.Order("topic.is_top desc, topic.is_essence desc,topic.see_count desc, topic.updated_at desc, topic.like_count desc,  topic.id desc")
	} else {
		query.Order("topic.is_top desc, topic.is_essence desc,topic.see_count desc, topic.updated_at desc, topic.like_count desc,  topic.id desc")
	}

	if by.Limit != 0 {
		query.Limit(by.Limit)
	}

	if by.Offset != 0 {
		query.Offset(by.Offset)
	}

	err := query.Find(&topList).Error
	if err != nil {
		return nil, 0, err
	}
	return topList, total, nil
}

func (d defaultTopicRepository) GetTopicPageList(by GetTopicPageListBy) (list []entity.Topic, total int64) {
	list = make([]entity.Topic, 0)

	query := d.ctx.DB.Model(&entity.Topic{})
	if by.ID > 0 {
		query.Where("topic.id = ?", by.ID)
	}
	if by.TopicTagId > 0 {
		query.Joins("inner join topic_tag on topic.id = topic_tag.topic_id").Where("topic_tag.tag_id = ?", by.TopicTagId)
	}
	if by.Status > 0 {
		query.Where("topic.status = ?", by.Status)
	}

	err := query.Count(&total).
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

func (d defaultTopicRepository) FindById(topicId int64) entity.Topic {
	var topic entity.Topic
	err := d.ctx.DB.Model(&entity.Topic{}).Preload("User").Preload("Tags").Where("id = ?", topicId).First(&topic).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return topic
}

func (d defaultTopicRepository) Save(topic *entity.Topic) error {
	return d.ctx.DB.Save(topic).Error
}

func (d defaultTopicRepository) AddTopicLikeCount(topicId int64, num int) error {
	db := d.ctx.DB.Model(entity.Topic{}).
		Where("id = ?", topicId)
	//避免点赞数为负数
	if num < 0 {
		db.Where("like_count >= ?", -num)
	}
	return db.Update("like_count", gorm.Expr("like_count + ?", num)).Error
}

func (d defaultTopicRepository) GetFlowPageList(by GetTopicFlowPageListBy) (list []entity.Topic, total int64) {
	list = make([]entity.Topic, 0)
	db := d.ctx.DB.Table(fmt.Sprintf("%s as flow", entity.TopicFlow{}.TableName())).
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

func (d defaultTopicRepository) UpdateColumn(id int64, key string, value interface{}) error {
	return d.ctx.DB.Model(&entity.Topic{}).Where("id = ?", id).Update(key, value).Error
}

func NewTopicModel(ctx *mioContext.MioContext) TopicModel {
	return defaultTopicRepository{
		ctx: ctx,
	}
}
