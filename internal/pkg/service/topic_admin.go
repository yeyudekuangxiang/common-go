package service

import (
	"bytes"
	"fmt"
	"github.com/medivhzhan/weapp/v3/security"
	"github.com/pkg/errors"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/pkg/wxoa"
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

func (srv TopicAdminService) GetTopicList(param repository.TopicListRequest) ([]entity.Topic, int64, error) {
	topList := make([]entity.Topic, 0)
	var total int64
	query := app.DB.Model(&entity.Topic{}).Preload("Tags").
		Joins("left join topic_tag on topic.id = topic_tag.topic_id")
	if param.TagId != 0 {
		query.Where("topic_tag.tag_id = ?", param.TagId)
	}
	if param.ID != 0 {
		query.Where("topic.id = ?", param.ID)
	}
	if param.Title != "" {
		query.Where("topic.title = ?", param.Title)
	}
	if param.UserId != 0 {
		query.Where("topic.user_id = >", param.UserId)
	}
	if param.UserName != "" {
		query.Where("topic.nickname = ?", param.UserName)
	}
	if param.Status > 0 {
		query.Where("topic.status = ?", param.Status)
	}
	if param.IsTop > 0 {
		query.Where("topic.s_top = ?", param.IsTop)
	}
	if param.IsEssence > 0 {
		query.Where("topic.is_essence = ?", param.IsEssence)
	}
	err := query.Count(&total).
		Group("topic.id").
		Order("is_top desc, is_essence desc, updated_at desc, created_at desc, id desc").
		Limit(param.Limit).
		Offset(param.Offset).
		Find(&topList).Error
	if err != nil {
		return topList, total, err
	}
	return topList, total, nil
}

//CreateTopic 创建文章
func (srv TopicAdminService) CreateTopic(userId int64, avatarUrl, nikeName, openid string, title, content string, tagIds []int64, images []string) error {
	//检查内容
	err := srv.checkMsgs(openid, content)
	if err != nil {
		return err
	}
	//处理images
	imageStr := strings.Join(images, ",")
	//tag
	tagModel := make([]entity.Tag, 0)
	for _, tagId := range tagIds {
		tagModel = append(tagModel, entity.Tag{
			Id: tagId,
		})
	}
	tag := srv.tag.GetById(tagIds[0])
	//topic
	topicModel := &entity.Topic{
		UserId:    userId,
		Title:     title,
		TopicTag:  tag.Name,
		Content:   content,
		ImageList: imageStr,
		Status:    entity.TopicStatusNeedVerify,
		Avatar:    avatarUrl,
		Nickname:  nikeName,
		Tags:      tagModel,
		CreatedAt: model.Time{},
		UpdatedAt: model.Time{},
	}
	if err := srv.topic.Save(topicModel); err != nil {
		return err
	}
	return nil
}

// UpdateTopic 更新帖子
func (srv TopicAdminService) UpdateTopic(openid string, topicId int64, title, content string, tagIds []int64, images []string) error {
	//检查内容
	err := srv.checkMsgs(openid, content)
	//查询记录是否存在
	topicModel := srv.topic.FindById(topicId)
	if topicModel.Id == 0 {
		return errors.New("该帖子不存在")
	}
	//处理images
	imageStr := strings.Join(images, ",")
	//tag
	tagModel := make([]entity.Tag, 0)
	for _, tagId := range tagIds {
		tagModel = append(tagModel, entity.Tag{
			Id: tagId,
		})
	}
	tag := srv.tag.GetById(tagIds[0])

	//更新帖子
	topicModel.Title = title
	topicModel.ImageList = imageStr
	topicModel.Content = content
	topicModel.TopicTag = tag.Name
	if err := app.DB.Model(&entity.Topic{}).Updates(topicModel).Error; err != nil {
		return err
	}
	err = app.DB.Model(&topicModel).Association("Tags").Replace(tagModel)
	if err != nil {
		return err
	}
	return nil
}

// DetailTopic 获取topic详情
func (srv TopicAdminService) DetailTopic(topicId int64) (entity.Topic, error) {
	//查询数据是否存在
	var topic entity.Topic
	app.DB.Model(&entity.Topic{}).Preload("Tags").Where("id = ?", topicId).Find(&topic)
	if topic.Id == 0 {
		return entity.Topic{}, errors.New("数据不存在")
	}
	return topic, nil
}

// DeleteTopic 删除（下架）
func (srv TopicAdminService) DeleteTopic(topicId int64) error {
	//查询数据是否存在
	var topic entity.Topic
	app.DB.Model(&entity.Topic{}).Preload("Tags").Where("id = ?", topicId).Find(&topic)
	if topic.Id == 0 {
		return errors.New("数据不存在")
	}
	err := app.DB.Model(&topic).Update("status", 4).Error
	if err != nil {
		return err
	}
	return nil
}

// Review 审核
func (srv TopicAdminService) Review(topicId int64, status int, content string) error {
	//查询数据是否存在
	var topic entity.Topic
	app.DB.Model(&topic).Where("id = ?", topicId).Find(&topic)
	if topic.Id == 0 {
		return errors.New("数据不存在")
	}
	if err := app.DB.Model(&topic).Updates(map[string]interface{}{"status": status, "content": content}).Error; err != nil {
		return err
	}
	//积分变动
	//pointService := NewPointService(context.NewMioContext())
	//point, err := pointService.IncUserPoint()
	//if err != nil {
	//	return err
	//}
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
	return nil
}

func (srv TopicAdminService) checkMsgs(openid, content string) error {
	length := len(content)
	if length > 2500 {
		s := []rune(content)
		var buffer bytes.Buffer
		for i, str := range s {
			buffer.WriteString(string(str))
			if i > 0 && (i+1)%2500 == 0 {
				params := &security.MsgSecCheckRequest{
					Content: "",
					Version: 2,
					Scene:   3,
					Openid:  openid,
				}
				err := srv.checkMsg(params)
				buffer.Reset()
				if err != nil {
					return err
				}
			}
		}
		params := &security.MsgSecCheckRequest{
			Content: "",
			Version: 2,
			Scene:   3,
			Openid:  openid,
		}
		err := srv.checkMsg(params)
		buffer.Reset()
		if err != nil {
			return err
		}
		return nil
	}
	//处理内容
	params := &security.MsgSecCheckRequest{
		Content: content,
		Version: 2,
		Scene:   3,
		Openid:  openid,
	}
	err := srv.checkMsg(params)
	if err != nil {
		return err
	}
	return nil
}

func (srv TopicAdminService) checkMsg(params *security.MsgSecCheckRequest) error {
	check, err := app.Weapp.NewSecurity().MsgSecCheck(params)
	if err != nil {
		return err
	}
	if check.ErrCode != 0 {
		return fmt.Errorf("check error: %s", check.ErrMSG)
	}
	return nil
}
