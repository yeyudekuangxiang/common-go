package service

import (
	"context"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"gorm.io/gorm"
	"mio/config"
	"mio/core/app"
	"mio/model/entity"
	"mio/repository"
	"time"
)

var DefaultTopicFlowService = TopicFlowService{repo: repository.DefaultTopicFlowRepository}

type TopicFlowService struct {
	repo repository.TopicFlowRepository
}

// CalculateSort 计算内容流排序 创建时间和更新时间未计算在内
func (t TopicFlowService) CalculateSort(topic entity.Topic, flow entity.TopicFlow) int {
	//有权重的 而且未看过
	if topic.Sort > 0 && flow.SeeCount == 0 {
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
	//我看过的 而且有权重的
	if flow.SeeCount > 0 && topic.Sort > 0 {
		return 6000 + topic.Sort
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

var initUserFlowPool, _ = ants.NewPool(50)

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
	topicList := make([]entity.Topic, 0)
	app.DB.FindInBatches(&topicList, 500, func(tx *gorm.DB, batch int) error {
		flowList := make([]entity.TopicFlow, 0)
		for _, topic := range topicList {
			flow := t.repo.FindBy(repository.FindTopicFlowBy{
				TopicId: topic.Id,
				UserId:  userId,
			})
			if flow.ID == 0 {
				flow = entity.TopicFlow{
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
func (t TopicFlowService) UpdateUserFlowSort(topic entity.Topic, flow *entity.TopicFlow) error {
	flow.Sort = t.CalculateSort(topic, *flow)
	return app.DB.Save(flow).Error
}

func (t TopicFlowService) AddUserFlowSeeCount(userId int64, topicId int64) {
	if userId == 0 || topicId == 0 {
		return
	}
	topic := DefaultTopicService.FindById(topicId)
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
	flow.Sort = t.CalculateSort(topic, flow)

	if err := t.repo.Save(&flow); err != nil {
		app.Logger.Error("更新查看数量失败", err)
	}
}
func (t TopicFlowService) AddUserFlowShowCount(userId int64, topicId int64) {
	if userId == 0 || topicId == 0 {
		return
	}
	topic := DefaultTopicService.FindById(topicId)
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
	flow.Sort = t.CalculateSort(topic, flow)

	if err := t.repo.Save(&flow); err != nil {
		app.Logger.Error("更新查看数量失败", err)
	}
}
