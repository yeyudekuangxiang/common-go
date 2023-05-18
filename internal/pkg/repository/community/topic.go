package community

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"mio/internal/pkg/core/app"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

type (
	TopicModel interface {
		GetTopicPageList(by repository.GetTopicPageListBy) (list []entity.Topic, total int64)
		FindById(topicId int64) *entity.Topic
		FindOneTopic(params repository.FindTopicParams) (*entity.Topic, error)
		FindOneTopicAndTag(params repository.FindTopicParams) (*entity.Topic, error)
		FindOneUnscoped(params repository.FindTopicParams) *entity.Topic
		Save(topic *entity.Topic) error
		AddTopicLikeCount(topicId int64, num int) error
		AddTopicSeeCount(topicId int64, num int) error
		GetFlowPageList(by repository.GetTopicFlowPageListBy) (list []entity.Topic, total int64)
		GetMyTopic(params repository.MyTopicListParams) ([]*entity.Topic, int64, error)
		GetTopicList(by repository.GetTopicPageListBy) ([]*entity.Topic, int64, error)
		ChangeTopicCollectionCount(id int64, column string, incr int) error
		Trans(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error
		GetTopicNotes(topicIds []int64) []*entity.Topic
		GetTopicListV2(by repository.GetTopicPageListBy) ([]*entity.Topic, error)
		GetTopList() ([]*entity.Topic, error)
		GetImportTopic() ([]*entity.Topic, error)
		SoftDelete(topic *entity.Topic) error
		Updates(topic *entity.Topic) error
		UpdateColumn(id int64, key string, value interface{}) error
		UpdatesColumn(cond repository.UpdatesTopicCond, upColumns map[string]interface{}) error
	}

	defaultTopicModel struct {
		ctx *mioContext.MioContext
	}
)

func (d defaultTopicModel) SoftDelete(topic *entity.Topic) error {
	if err := d.ctx.DB.Delete(topic).Error; err != nil {
		return err
	}
	return nil
}

func (d defaultTopicModel) Updates(topic *entity.Topic) error {
	db := d.ctx.DB
	if topic.Type == TopicTypeActivity {
		activity := topic.Activity
		err := db.Save(&activity).Error
		if err != nil {
			return err
		}
	}
	err := db.Omit(clause.Associations).Save(topic).Error
	if err != nil {
		return err
	}

	return err
}

func (d defaultTopicModel) GetImportTopic() ([]*entity.Topic, error) {
	var resp []*entity.Topic
	err := d.ctx.DB.Model(&entity.Topic{}).
		Where("import_id != 0").
		Where("is_essence = 0").
		Select("id,user_id").
		Group("id,user_id").
		Order("id asc").
		Find(&resp).Error
	if err != nil {
		return []*entity.Topic{}, err
	}
	return resp, nil
}

func (d defaultTopicModel) GetTopList() ([]*entity.Topic, error) {
	topList := make([]*entity.Topic, 0)

	query := d.ctx.DB.Model(&entity.Topic{}).
		Preload("User").
		Preload("Tags").
		Preload("Activity").
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

	err := query.Where("topic.is_top = 1").Where("status = 3").Group("topic.id").Find(&topList).Error

	if err != nil {
		return nil, err
	}
	return topList, nil
}

func (d defaultTopicModel) FindOneTopic(params repository.FindTopicParams) (*entity.Topic, error) {
	var resp entity.Topic
	query := d.ctx.DB.Model(&entity.Topic{}).Preload("User")
	if params.TopicId != 0 {
		query = query.Where("id = ?", params.TopicId)
	}
	if params.UserId != 0 {
		query = query.Where("user_id = ?", params.UserId)
	}
	if params.Type != nil {
		query = query.Where("type = ?", *params.Type)
	}
	if params.Status != 0 {
		query = query.Where("status = ?", params.Status)
	}
	err := query.First(&resp).Error
	switch err {
	case nil:
		return &resp, nil
	case gorm.ErrRecordNotFound:
		return nil, entity.ErrNotFount
	default:
		return nil, err
	}
}

func (d defaultTopicModel) GetTopicNotes(topicIds []int64) []*entity.Topic {
	topList := make([]*entity.Topic, 0)
	err := d.ctx.DB.Model(&entity.Topic{}).
		Preload("User").
		Preload("Tags").
		Preload("Activity").
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

func (d defaultTopicModel) Trans(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return d.ctx.DB.Transaction(fc, opts...)
}

func (d defaultTopicModel) ChangeTopicCollectionCount(id int64, column string, incr int) error {
	return d.ctx.DB.Model(&entity.Topic{}).Where("id = ?", id).Update(column, gorm.Expr(column+"+?", incr)).Error
}

func (d defaultTopicModel) GetMyTopic(by repository.MyTopicListParams) ([]*entity.Topic, int64, error) {
	topList := make([]*entity.Topic, 0)
	var total int64
	query := d.ctx.DB.Model(&entity.Topic{}).
		Preload("User").
		Preload("Tags").
		Preload("Activity").
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

	if by.Type >= 0 {
		query.Where("topic.type = ?", by.Type)
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

func (d defaultTopicModel) GetTopicList(params repository.GetTopicPageListBy) ([]*entity.Topic, int64, error) {
	topList := make([]*entity.Topic, 0)
	var total int64
	query := d.ctx.DB.Model(&entity.Topic{}).
		Preload("User").
		Preload("Tags").
		Preload("Activity").
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

	if params.ID != 0 {
		query.Where("topic.id = ?", params.ID)
	} else if len(params.Ids) > 0 {
		query.Where("topic.id in ?", params.Ids)
	}

	if params.Label == "recommend" && len(params.Rids) > 0 {
		query.Where("topic.id in ?", params.Rids)
	}

	if params.UserId != 0 {
		query.Where("topic.user_id = ?", params.UserId)
	}

	if params.Label == "activity" {
		query.Where("topic.type = ?", 1)
	}

	if params.Status != 0 {
		query.Where("topic.status = ?", params.Status)
	} else {
		query.Where("topic.status = ?", 3)
	}

	if params.TopicTagId != 0 {
		query.Joins("inner join topic_tag on topic.id = topic_tag.topic_id").Where("topic_tag.tag_id = ?", params.TopicTagId)
	}

	query = query.Count(&total).
		Group("topic.id").
		Order("id desc")

	if params.Limit != 0 {
		query.Limit(params.Limit)
	}

	if params.Offset != 0 {
		query.Offset(params.Offset)
	}

	err := query.Find(&topList).Error
	if err != nil {
		return nil, 0, err
	}
	return topList, total, nil
}

func (d defaultTopicModel) GetTopicListV2(by repository.GetTopicPageListBy) ([]*entity.Topic, error) {
	topList := make([]*entity.Topic, 0)

	query := d.ctx.DB.Model(&entity.Topic{}).
		Preload("User").
		Preload("Tags").
		Preload("Activity").
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

	query.Where("topic.id in ?", by.Rids).Where("status = 3")

	err := query.Group("topic.id").Find(&topList).Error

	if err != nil {
		return nil, err
	}
	return topList, nil
}

func (d defaultTopicModel) GetTopicPageList(by repository.GetTopicPageListBy) (list []entity.Topic, total int64) {
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

func (d defaultTopicModel) FindById(topicId int64) *entity.Topic {
	var resp entity.Topic
	err := d.ctx.DB.Model(&entity.Topic{}).
		Preload("User").
		Preload("Tags").
		Preload("Activity").
		Where("id = ?", topicId).
		Where("status = ?", 3).
		First(&resp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	return &resp
}

func (d defaultTopicModel) FindOneUnscoped(params repository.FindTopicParams) *entity.Topic {
	var resp entity.Topic
	query := d.ctx.DB.Model(&entity.Topic{}).
		Preload("User").
		Preload("Tags").
		Preload("Activity").
		Unscoped()
	if params.TopicId != 0 {
		query.Where("id = ?", params.TopicId)
	}
	if params.Status != 0 {
		query.Where("status = ?", params.Status)
	}
	if params.Type != nil {
		query.Where("type = ?", *params.Type)
	}
	if params.UserId != 0 {
		query.Where("user_id = ?", params.UserId)
	}
	err := query.First(&resp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	return &resp
}

func (d defaultTopicModel) Save(topic *entity.Topic) error {
	return d.ctx.DB.Save(topic).Error
}

func (d defaultTopicModel) AddTopicLikeCount(topicId int64, num int) error {
	db := d.ctx.DB.Model(&entity.Topic{}).
		Where("id = ?", topicId)
	//避免点赞数为负数
	if num < 0 {
		db.Where("like_count >= ?", -num)
	}
	return db.Update("like_count", gorm.Expr("like_count + ?", num)).Error
}

func (d defaultTopicModel) AddTopicSeeCount(topicId int64, num int) error {
	db := d.ctx.DB.Model(&entity.Topic{}).
		Where("id = ?", topicId)
	//避免点赞数为负数
	if num < 0 {
		db.Where("see_count >= ?", -num)
	}
	return db.Update("see_count", gorm.Expr("see_count + ?", num)).Error
}

func (d defaultTopicModel) GetFlowPageList(by repository.GetTopicFlowPageListBy) (list []entity.Topic, total int64) {
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

func (d defaultTopicModel) UpdateColumn(id int64, key string, value interface{}) error {
	return d.ctx.DB.Model(&entity.Topic{}).Where("id = ?", id).Update(key, value).Error
}

func (d defaultTopicModel) UpdatesColumn(cond repository.UpdatesTopicCond, upColumns map[string]interface{}) error {
	query := d.ctx.DB.Model(&entity.Topic{})
	if cond.Id != 0 {
		query.Where("id = ?", cond.Id)
	} else if len(cond.Ids) > 0 {
		query.Where("id in (?)", cond.Ids)
	}

	if cond.IsEssence != 0 {
		if cond.IsEssence < 0 {
			query.Where("is_essence = ?", 0)
		} else {
			query.Where("is_essence = ?", 1)
		}
	}

	if cond.IsTop != 0 {
		if cond.IsTop < 0 {
			query.Where("is_top = ?", 0)
		} else {
			query.Where("is_top = ?", 1)
		}
	}

	err := query.Updates(upColumns).Error
	if err != nil {
		return err
	}

	return nil
}

func (d defaultTopicModel) FindOneTopicAndTag(params repository.FindTopicParams) (*entity.Topic, error) {
	var resp entity.Topic
	query := d.ctx.DB.Model(&entity.Topic{}).Preload("User").Preload("Tags")
	if params.TopicId != 0 {
		query = query.Where("id = ?", params.TopicId)
	}
	if params.UserId != 0 {
		query = query.Where("user_id = ?", params.UserId)
	}
	if params.Type != nil {
		query = query.Where("type = ?", *params.Type)
	}
	if params.Status != 0 {
		query = query.Where("status = ?", params.Status)
	}
	err := query.First(&resp).Error
	switch err {
	case nil:
		return &resp, nil
	case gorm.ErrRecordNotFound:
		return nil, entity.ErrNotFount
	default:
		return nil, err
	}
}

func NewTopicModel(ctx *mioContext.MioContext) TopicModel {
	return defaultTopicModel{
		ctx: ctx,
	}
}
