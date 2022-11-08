package service

import (
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/pkg/errno"
	"mio/pkg/wxoa"
	"strconv"
	"strings"
)

func NewTopicAdminService(ctx *context.MioContext) TopicAdminService {
	return TopicAdminService{
		ctx:        ctx,
		topicModel: repository.NewTopicModel(ctx),
		tag:        repository.NewTagModel(ctx),
	}
}

type TopicAdminService struct {
	ctx         *context.MioContext
	topicModel  repository.TopicModel
	tag         repository.TagModel
	TokenServer *wxoa.AccessTokenServer
}

func (srv TopicAdminService) GetTopicList(param repository.TopicListRequest) ([]*entity.Topic, int64, error) {
	topList := make([]*entity.Topic, 0)
	var total int64
	query := app.DB.Model(&entity.Topic{})
	if param.ID != 0 {
		query.Where("topic.id = ?", param.ID)
	}

	if param.Title != "" {
		query.Where("topic.title like ?", "%"+param.Title+"%")
	}

	if param.UserId != 0 {
		query.Where("topic.user_id = ?", param.UserId)
	}

	if param.Status > 0 {
		query.Where("topic.status = ?", param.Status)
	}

	if param.IsTop > 0 {
		query.Where("topic.is_top = ?", param.IsTop)
	}

	if param.IsEssence > 0 {
		query.Where("topic.is_essence = ?", param.IsEssence)
	}

	if param.TagId != 0 {
		query.Joins("left join topic_tag on topic.id = topic_tag.topic_id").Where("topic_tag.tag_id = ?", param.TagId)
	} else if len(param.TagIds) > 0 {
		query.Joins("left join topic_tag on topic.id = topic_tag.topic_id").Where("topic_tag.tag_id in (?)", param.TagIds)
	}

	if param.UserName != "" || param.IsPartners != 0 || param.Position != "" {
		query.Joins("left join \"user\" on \"user\".id = topic.user_id")
		if param.UserName != "" {
			query.Where("\"user\".nick_name = ?", param.UserName)
		}

		if param.IsPartners != 0 {
			query.Where("\"user\".partners = ?", param.IsPartners)
		}

		if param.Position != "" {
			query.Where("\"user\".position = ?", param.Position)
		}

	}

	query.Preload("Tags").Preload("User")
	err := query.
		Count(&total).
		Group("topic.id").
		Order("topic.id desc").
		Limit(param.Limit).
		Offset(param.Offset).
		Find(&topList).Error
	if err != nil {
		return topList, total, err
	}
	return topList, total, nil
}

//CreateTopic 创建文章
func (srv TopicAdminService) CreateTopic(userId int64, title, content string, tagIds []int64, images []string) error {
	//处理images
	imageStr := strings.Join(images, ",")
	//topic
	topicModel := &entity.Topic{
		UserId:    userId,
		Title:     title,
		Content:   content,
		ImageList: imageStr,
		Status:    entity.TopicStatusNeedVerify,
		CreatedAt: model.Time{},
		UpdatedAt: model.Time{},
	}

	//tag
	if len(tagIds) > 0 {
		tagModel := make([]entity.Tag, 0)
		for _, tagId := range tagIds {
			tagModel = append(tagModel, entity.Tag{
				Id: tagId,
			})
		}
		tag := srv.tag.GetById(tagIds[0])
		topicModel.Tags = tagModel
		topicModel.TopicTag = tag.Name
		topicModel.TopicTagId = strconv.FormatInt(tag.Id, 10)
	}

	if err := srv.topicModel.Save(topicModel); err != nil {
		return err
	}
	return nil
}

// UpdateTopic 更新帖子
func (srv TopicAdminService) UpdateTopic(topicId int64, title, content string, tagIds []int64, images []string) error {
	//查询记录是否存在
	topicModel := srv.topicModel.FindById(topicId)
	if topicModel.Id == 0 {
		return errno.ErrCommon.WithMessage("该帖子不存在")
	}
	//处理images
	imageStr := strings.Join(images, ",")

	//更新帖子
	topicModel.Title = title
	topicModel.ImageList = imageStr
	topicModel.Content = content

	//tag
	if len(tagIds) > 0 {
		tagModel := make([]entity.Tag, 0)
		for _, tagId := range tagIds {
			tagModel = append(tagModel, entity.Tag{
				Id: tagId,
			})
		}
		tag := srv.tag.GetById(tagIds[0])
		topicModel.TopicTag = tag.Name
		topicModel.TopicTagId = strconv.FormatInt(tag.Id, 10)
		err := app.DB.Model(&topicModel).Association("Tags").Replace(tagModel)
		if err != nil {
			return err
		}
	} else {
		topicModel.TopicTag = ""
		topicModel.TopicTagId = ""
		err := app.DB.Model(&topicModel).Association("Tags").Clear()
		if err != nil {
			return err
		}
	}

	if err := app.DB.Model(&topicModel).Updates(topicModel).Error; err != nil {
		return err
	}
	return nil
}

// DetailTopic 获取topic详情
func (srv TopicAdminService) DetailTopic(topicId int64) (entity.Topic, error) {
	//查询数据是否存在
	var topic entity.Topic
	app.DB.Model(&entity.Topic{}).Preload("Tags").Preload("User").Where("id = ?", topicId).Find(&topic)
	if topic.Id == 0 {
		return entity.Topic{}, errno.ErrCommon.WithMessage("数据不存在")
	}
	return topic, nil
}

// DeleteTopic 删除（下架）
func (srv TopicAdminService) DeleteTopic(topicId int64, reason string) (*entity.Topic, error) {
	//查询数据是否存在
	topic := srv.topicModel.FindById(topicId)
	if topic.Id == 0 {
		return nil, errno.ErrCommon.WithMessage("数据不存在")
	}

	if topic.Status == 4 {
		return nil, nil
	}

	topic.Status = 4
	topic.DelReason = reason

	err := srv.topicModel.Save(topic)
	if err != nil {
		return nil, err
	}

	return topic, nil
}

// Review 审核
func (srv TopicAdminService) Review(topicId int64, status int, reason string) (entity.Topic, bool, error) {
	//查询数据是否存在
	topic := srv.topicModel.FindById(topicId)
	var isFirst bool

	if topic.Id == 0 {
		return entity.Topic{}, isFirst, errno.ErrCommon.WithMessage("数据不存在")
	}

	if status == entity.TopicStatusPublished {
		if topic.PushTime.IsZero() {
			isFirst = true
		}
		topic.Status = entity.TopicStatusPublished
		topic.PushTime = model.NewTime()
	}

	if status == entity.TopicStatusHidden {
		if topic.DownTime.IsZero() {
			isFirst = true
		}
		topic.Status = entity.TopicStatusHidden
		topic.DownTime = model.NewTime()
		topic.DelReason = reason
	}

	if status == entity.TopicStatusVerifyFailed {
		topic.Status = entity.TopicStatusVerifyFailed
	}

	//更新帖子
	err := srv.topicModel.Save(topic)

	if err != nil {
		return entity.Topic{}, isFirst, err
	}

	return *topic, isFirst, nil
}

// Top 置顶
func (srv TopicAdminService) Top(topicId int64, isTop int) (*entity.Topic, bool, error) {
	//查询数据是否存在
	topic := srv.topicModel.FindById(topicId)

	var isFirst bool

	if topic.Id == 0 {
		return &entity.Topic{}, isFirst, errno.ErrCommon.WithMessage("数据不存在")
	}
	update := entity.Topic{IsTop: isTop}
	if isTop == 1 {
		if topic.TopTime.IsZero() {
			isFirst = true
		}
		update.TopTime = model.NewTime()
	}

	if err := app.DB.Model(&topic).Updates(update).Error; err != nil {
		return &entity.Topic{}, isFirst, err
	}

	return topic, isFirst, nil
}

// Essence 精华
func (srv TopicAdminService) Essence(topicId int64, isEssence int) (*entity.Topic, bool, error) {
	//查询数据是否存在
	topic := srv.topicModel.FindById(topicId)
	var isFirst bool
	if topic.Id == 0 {
		return &entity.Topic{}, false, errno.ErrCommon.WithMessage("数据不存在")
	}

	update := entity.Topic{
		IsEssence: isEssence,
	}

	if isEssence == 1 {
		if topic.EssenceTime.IsZero() {
			isFirst = true
		}
		update.EssenceTime = model.NewTime()
	}

	if err := app.DB.Model(&topic).Updates(update).Error; err != nil {
		return &entity.Topic{}, false, err
	}

	return topic, isFirst, nil
}

func (srv TopicAdminService) GetCommentCount(ids []int64) (result []CommentCount) {
	app.DB.Model(&entity.CommentIndex{}).Select("obj_id as topic_id, count(*) as total").
		Where("obj_id in ?", ids).
		Group("obj_id").
		Find(&result)
	return result
}
