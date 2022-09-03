package service

import (
	"github.com/pkg/errors"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/pkg/wxoa"
	"strconv"
	"strings"
)

var DefaultTopicAdminService = NewTopicAdminService()

func NewTopicAdminService() TopicAdminService {
	return TopicAdminService{
		topic: repository.DefaultTopicRepository,
		tag:   repository.DefaultTagRepository,
	}
}

type TopicAdminService struct {
	topic       repository.ITopicRepository
	tag         repository.ITagRepository
	TokenServer *wxoa.AccessTokenServer
}

func (srv TopicAdminService) GetTopicList(param repository.TopicListRequest) ([]*entity.Topic, int64, error) {
	topList := make([]*entity.Topic, 0)
	var total int64
	query := app.DB.Model(&entity.Topic{}).Preload("Tags").Preload("User")

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
	}
	err := query.Group("topic.id").
		Order("is_top desc, is_essence desc, updated_at desc, created_at desc, id desc").
		Limit(param.Limit).
		Offset(param.Offset).
		Find(&topList).Error
	if err != nil {
		return topList, total, err
	}
	total = int64(len(topList))
	return topList, total, nil
}

//CreateTopic 创建文章
func (srv TopicAdminService) CreateTopic(userId int64, avatarUrl, nikeName, title, content string, tagIds []int64, images []string) error {
	//检查内容
	//err := srv.checkMsgs(openid, content)
	//if err != nil {
	//	return err
	//}
	//处理images
	imageStr := strings.Join(images, ",")
	//topic
	topicModel := &entity.Topic{
		UserId:    userId,
		Title:     title,
		Content:   content,
		ImageList: imageStr,
		Status:    entity.TopicStatusNeedVerify,
		Avatar:    avatarUrl,
		Nickname:  "绿喵",
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

	if err := srv.topic.Save(topicModel); err != nil {
		return err
	}
	return nil
}

// UpdateTopic 更新帖子
func (srv TopicAdminService) UpdateTopic(topicId int64, title, content string, tagIds []int64, images []string) error {
	//检查内容
	//err := srv.checkMsgs(openid, content)
	//查询记录是否存在
	topicModel := srv.topic.FindById(topicId)
	if topicModel.Id == 0 {
		return errors.New("该帖子不存在")
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
		return entity.Topic{}, errors.New("数据不存在")
	}
	return topic, nil
}

// DeleteTopic 删除（下架）
func (srv TopicAdminService) DeleteTopic(topicId int64, reason string) error {
	//查询数据是否存在
	var topic entity.Topic
	app.DB.Model(&entity.Topic{}).Preload("Tags").Where("id = ?", topicId).Find(&topic)
	if topic.Id == 0 {
		return errors.New("数据不存在")
	}
	err := app.DB.Model(&topic).Updates(entity.Topic{Status: 4, DelReason: reason}).Error
	if err != nil {
		return err
	}
	return nil
}

// Review 审核
func (srv TopicAdminService) Review(topicId int64, status int, reason string) error {
	//查询数据是否存在
	var topic entity.Topic
	app.DB.Model(&topic).Where("id = ?", topicId).Find(&topic)
	if topic.Id == 0 {
		return errors.New("数据不存在")
	}

	if err := app.DB.Model(&topic).Updates(entity.Topic{Status: entity.TopicStatus(status), DelReason: reason}).Error; err != nil {
		return err
	}
	ctx := context.NewMioContext()

	pointService := NewPointService(ctx)
	pointTransService := NewPointTransactionService(ctx)
	user, _ := DefaultUserService.GetUserById(topic.UserId)

	by, err := pointTransService.FindBy(repository.FindPointTransactionBy{
		OpenId:       user.OpenId,
		Type:         string(entity.POINT_ARTICLE),
		AdditionInfo: strconv.FormatInt(topic.Id, 10),
	})
	if err != nil {
		return err
	}

	var point int64
	if topic.Status == 3 && status == 4 {
		point = -int64(entity.PointCollectValueMap[entity.POINT_ARTICLE])
	}
	title := topic.Title
	if len([]rune(title)) > 8 {
		title = string([]rune(title)[0:8]) + "..."
	}
	//发放积分
	if status == 3 && by.ID == 0 {
		point = int64(entity.PointCollectValueMap[entity.POINT_ARTICLE])
		_, _ = pointService.IncUserPoint(srv_types.IncUserPointDTO{
			OpenId:       user.OpenId,
			Type:         entity.POINT_ARTICLE,
			BizId:        util.UUID(),
			ChangePoint:  point,
			AdminId:      0,
			Note:         "笔记 " + title + "... 审核通过，发布成功",
			AdditionInfo: strconv.FormatInt(topic.Id, 10),
		})
	}
	return nil
}

// Top 置顶
func (srv TopicAdminService) Top(topicId int64, isTop int) error {
	//查询数据是否存在
	var topic entity.Topic
	app.DB.Model(&topic).Where("id = ?", topicId).Find(&topic)
	if topic.Id == 0 {
		return errors.New("数据不存在")
	}
	if err := app.DB.Model(&topic).Update("is_top", isTop).Error; err != nil {
		return err
	}
	return nil
}

// Essence 精华
func (srv TopicAdminService) Essence(topicId int64, isEssence int) error {
	//查询数据是否存在
	var topic entity.Topic
	app.DB.Model(&topic).Where("id = ?", topicId).Find(&topic)
	if topic.Id == 0 {
		return errors.New("数据不存在")
	}
	if err := app.DB.Model(&topic).Update("is_essence", isEssence).Error; err != nil {
		return err
	}
	//发放积分
	if isEssence == 1 {
		title := topic.Title
		if len([]rune(title)) > 8 {
			title = string([]rune(title)[0:8]) + "..."
		}
		user, _ := DefaultUserService.GetUserById(topic.UserId)
		pointService := NewPointService(context.NewMioContext())
		_, _ = pointService.IncUserPoint(srv_types.IncUserPointDTO{
			OpenId:       user.OpenId,
			Type:         entity.POINT_RECOMMEND,
			BizId:        util.UUID(),
			ChangePoint:  int64(entity.PointCollectValueMap[entity.POINT_RECOMMEND]),
			AdminId:      0,
			Note:         "笔记 \"" + title + "...\" 被设为精华",
			AdditionInfo: strconv.FormatInt(topic.Id, 10),
		})
	}
	return nil
}
