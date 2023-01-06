package community

import (
	"context"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"gorm.io/gorm"
	"mio/config"
	"mio/internal/pkg/core/app"
	context2 "mio/internal/pkg/core/context"
	entity2 "mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"time"
)

var initUserFlowPool, _ = ants.NewPool(50)

var DefaultTopicFlowService = TopicFlowService{
	repo:       repository.DefaultTopicFlowRepository,
	topicModel: repository.NewTopicModel(context2.NewMioContext()),
}

type TopicFlowService struct {
	repo       repository.TopicFlowRepository
	topicModel repository.TopicModel
}

// CalculateSort 计算内容流排序 创建时间和更新时间未计算在内
func (t TopicFlowService) CalculateSort(topic entity2.Topic, flow entity2.TopicFlow) int {
	//有权重的
	if topic.Sort > 0 {
		return 9000 + topic.Sort
	}
	//所有人未看过 而且对我曝光次数小于5次
	if flow.ShowCount < 5 && topic.SeeCount == 0 {
		return 8000
	}
	//对我曝光次数小于5次 而且未看过
	if flow.ShowCount < 5 && flow.SeeCount == 0 {
		return 7000
	}
	//我看过的
	if flow.SeeCount > 0 {
		return 6000
	}
	//对我曝光次数大于5次 答案是未看过的
	if flow.ShowCount >= 5 && flow.SeeCount == 0 {
		return 5000
	}
	return 0
}

// InitUserFlow 同步用户内容流
func (t TopicFlowService) InitUserFlow(userId int64) {
	if userId == 0 {
		return
	}
	//防止多个多个应用副本同时初始化
	limitKey := fmt.Sprintf(config.RedisKey.InitTopicFlowLimit, userId)
	ok, err := app.Redis.SetNX(context.Background(), limitKey, 1, time.Hour*24).Result()
	if err != nil && !ok {
		app.Logger.Info("limit", ok, err)
		return
	}
	app.Logger.Info("limit", ok, err)
	app.Logger.Error(`开始同步内容流`, userId)
	topicList := make([]entity2.Topic, 0)
	app.DB.FindInBatches(&topicList, 500, func(tx *gorm.DB, batch int) error {
		flowList := make([]entity2.TopicFlow, 0)
		for _, topic := range topicList {
			flow := t.repo.FindBy(repository.FindTopicFlowBy{
				TopicId: topic.Id,
				UserId:  userId,
			})
			if flow.ID == 0 {
				flow = entity2.TopicFlow{
					UserId:         userId,
					TopicId:        topic.Id,
					SeeCount:       0,
					ShowCount:      0,
					TopicCreatedAt: topic.CreatedAt,
					TopicUpdatedAt: topic.UpdatedAt,
				}
			}
			flow.Sort = t.CalculateSort(topic, flow)
			flowList = append(flowList, flow)
		}
		err := app.DB.Transaction(func(tx *gorm.DB) error {
			return tx.Save(&flowList).Error
		})
		if err != nil {
			app.Logger.Error(`同步内容流失败`, userId, batch, 500, err)
		} else {
			app.Logger.Info(`同步内容流成功`, userId, batch, 500, err)
		}
		return nil
	})
}

// InitUserFlowByMq 同步用户内容流 前期用go携程同步 后期用消息队列同步
func (t TopicFlowService) InitUserFlowByMq(userId int64) {
	err := initUserFlowPool.Submit(func() {
		t.InitUserFlow(userId)
	})
	if err != nil {
		app.Logger.Error("提交初始化内容流任务失败", userId, err)
	}
}

// UpdateUserFlowSort 更新内容权重后重新计算内容流排序
func (t TopicFlowService) UpdateUserFlowSort(topic entity2.Topic, flow *entity2.TopicFlow) error {
	flow.Sort = t.CalculateSort(topic, *flow)
	return app.DB.Save(flow).Error
}

// AddUserFlowSeeCount 增加用户内容流查看次数
func (t TopicFlowService) AddUserFlowSeeCount(userId int64, topicId int64) {
	if userId == 0 || topicId == 0 {
		return
	}

	topic := t.topicModel.FindById(topicId)
	if topic.Id == 0 {
		return
	}
	flow := t.repo.FindBy(repository.FindTopicFlowBy{
		TopicId: topicId,
		UserId:  userId,
	})
	if flow.ID == 0 {
		return
	}
	flow.SeeCount++
	flow.Sort = t.CalculateSort(*topic, flow)

	if err := t.repo.Save(&flow); err != nil {
		app.Logger.Error("更新查看数量失败", err)
	}
}

// AddUserFlowShowCount 增加用户内容流曝光次数
func (t TopicFlowService) AddUserFlowShowCount(userId int64, topicId int64) {
	if userId == 0 || topicId == 0 {
		return
	}
	topic := t.topicModel.FindById(topicId)
	if topic.Id == 0 {
		return
	}
	flow := t.repo.FindBy(repository.FindTopicFlowBy{
		TopicId: topicId,
		UserId:  userId,
	})
	if flow.ID == 0 {
		return
	}
	flow.ShowCount++
	flow.Sort = t.CalculateSort(*topic, flow)

	if err := t.repo.Save(&flow); err != nil {
		app.Logger.Error("更新查看数量失败", err)
	}
}

// AfterUpdateTopic 新增topic 、更新topic权重、更新topic查看次数 后调用此方法 重新计算内容流排序
func (t TopicFlowService) AfterUpdateTopic(topicId int64) {
	topic := t.topicModel.FindById(topicId)
	if topic.Id == 0 {
		return
	}
	type UserIds struct {
		UserId int64
	}
	userIds := make([]UserIds, 0)
	app.DB.Table("(?) as t", app.DB.Model(entity2.TopicFlow{}).Select("user_id").Group("user_id")).FindInBatches(&userIds, 200, func(tx *gorm.DB, batch int) error {
		for _, item := range userIds {
			topicFlow := entity2.TopicFlow{}
			app.DB.Where("user_id = ? and topic_id = ?", item.UserId, topicId).First(&topicFlow)
			if topicFlow.ID != 0 {
				topicFlow.Sort = t.CalculateSort(*topic, topicFlow)
				if err := app.DB.Save(&topicFlow).Error; err != nil {
					app.Logger.Error("同步topic_flow失败", topicId, item.UserId, err)
				}
			} else {
				topicFlow = entity2.TopicFlow{
					UserId:         item.UserId,
					TopicId:        topicId,
					SeeCount:       0,
					ShowCount:      0,
					TopicCreatedAt: topic.CreatedAt,
					TopicUpdatedAt: topic.UpdatedAt,
				}
				topicFlow.Sort = t.CalculateSort(*topic, topicFlow)
				if err := app.DB.Create(&topicFlow).Error; err != nil {
					app.Logger.Error("同步topic_flow失败", topicId, item.UserId, err)
				}
			}
		}
		return nil
	})
}

// AfterUpdateTopicByMq 新增topic 、更新topic权重、更新topic查看次数 后调用此方法 重新计算内容流排序
func (t TopicFlowService) AfterUpdateTopicByMq(topicId int64) {
	err := initUserFlowPool.Submit(func() {
		t.AfterUpdateTopic(topicId)
	})
	if err != nil {
		app.Logger.Error("提交AfterUpdateTopicByMq任务失败", topicId, err)
	}
}
