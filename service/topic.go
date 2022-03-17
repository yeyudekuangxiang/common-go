package service

import (
	"fmt"
	"github.com/medivhzhan/weapp/v3"
	"github.com/pkg/errors"
	"mio/core/app"
	"mio/internal/wxapp"
	"mio/model/entity"
	"mio/repository"
)

var DefaultTopicService = NewTopicService(repository.DefaultTopicRepository)

func NewTopicService(r repository.ITopicRepository) TopicService {
	return TopicService{
		r: r,
	}
}

type TopicService struct {
	r repository.ITopicRepository
}

func (u TopicService) fillTopicList(topicList []entity.Topic, userId int64) ([]TopicDetail, error) {
	//查询点赞信息
	topicIds := make([]int64, 0)
	for _, topic := range topicList {
		topicIds = append(topicIds, topic.Id)
	}
	topicLikeMap := make(map[int64]bool)
	if userId > 0 {
		likeList := repository.TopicLikeRepository{DB: app.DB}.GetListBy(repository.GetTopicLikeListBy{
			TopicIds: topicIds,
			UserId:   userId,
		})
		for _, like := range likeList {
			topicLikeMap[int64(like.TopicId)] = like.Status == 1
		}
	}

	//整理数据
	detailList := make([]TopicDetail, 0)
	for _, topic := range topicList {
		detailList = append(detailList, TopicDetail{
			Topic:         topic,
			IsLike:        topicLikeMap[topic.Id],
			UpdatedAtDate: topic.UpdatedAt.Format("01-02"),
		})
	}

	return detailList, nil
}
func (u TopicService) GetTopicDetailPageList(param repository.GetTopicPageListBy) ([]TopicDetail, int64, error) {
	list, total := u.r.GetTopicPageList(param)

	//更新曝光和查看次数
	u.UpdateTopicListShowCount(list, param.UserId)
	if param.ID != 0 && len(list) > 0 {
		app.Logger.Info("更新查看次数", list[0].Id, param.UserId)
		u.UpdateTopicListSeeCount(list[0].Id, param.UserId)
	}

	detailList, err := u.fillTopicList(list, param.UserId)
	if err != nil {
		return nil, 0, err
	}
	return detailList, total, nil
}

// GetTopicFlowPageList1 获取用户内容流 如果没有数据则从topic表返回 同时进行初始化操作
/*func (u TopicService) GetTopicFlowPageList1(param repository.GetTopicPageListBy) ([]TopicDetail, int64, error) {
	flowList, total, err := DefaultTopicFlowService.GetPageList(repository.GetTopicFlowPageListBy{
		UserId:     param.UserId,
		Offset:     param.Offset,
		Limit:      param.Limit,
		TopicId:    param.ID,
		TopicTagId: param.TopicTagId,
	})
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		DefaultTopicFlowService.InitUserFlowByMq(param.UserId)
		return u.GetTopicDetailPageList(param)
	}

	topicIds := make([]int64, 0)
	for _, topicFlow := range flowList {
		topicIds = append(topicIds, topicFlow.TopicId)
	}

	topicList := u.r.GetTopicList(repository.GetTopicListBy{
		TopicIds: topicIds,
	})
	if err != nil {
		return nil, 0, err
	}

	//更新曝光和查看次数
	u.UpdateTopicListShowCount(topicList, param.UserId)
	if param.ID != 0 && len(topicList) > 0 {
		u.UpdateTopicListSeeCount(topicList[0].Id, param.UserId)
	}

	topicList = u.sortTopicListByIds(topicList, topicIds)

	topicDetailList, err := u.fillTopicList(topicList, param.UserId)

	if err != nil {
		return nil, 0, err
	}
	return topicDetailList, total, nil
}*/
// GetTopicFlowPageList 获取用户内容流 如果没有数据则从topic表返回 同时进行初始化操作
func (u TopicService) GetTopicFlowPageList(param repository.GetTopicPageListBy) ([]TopicDetail, int64, error) {

	topicList, total, err := u.r.GetFlowPageList(repository.GetTopicFlowPageListBy{
		Offset:     param.Offset,
		Limit:      param.Limit,
		UserId:     param.UserId,
		TopicId:    param.ID,
		TopicTagId: param.TopicTagId,
	})
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		DefaultTopicFlowService.InitUserFlowByMq(param.UserId)
		return u.GetTopicDetailPageList(param)
	}

	//更新曝光和查看次数
	u.UpdateTopicListShowCount(topicList, param.UserId)

	if param.ID != 0 && len(topicList) > 0 {
		app.Logger.Info("更新查看次数", param.UserId, topicList[0].Id)
		u.UpdateTopicListSeeCount(topicList[0].Id, param.UserId)
	}

	topicDetailList, err := u.fillTopicList(topicList, param.UserId)

	if err != nil {
		return nil, 0, err
	}
	return topicDetailList, total, nil
}
func (u TopicService) UpdateTopicListSeeCount(topicId int64, userId int64) {
	err := initUserFlowPool.Submit(func() {
		topic := u.r.FindById(topicId)
		if topic.Id == 0 {
			return
		}
		topic.SeeCount++
		if err := u.r.Save(&topic); err != nil {
			app.Logger.Error("更新topic查看次数失败", topicId, userId)
			return
		}
		DefaultTopicFlowService.AddUserFlowSeeCount(userId, topicId)
	})
	if err != nil {
		app.Logger.Error("提交更新topic查看次数任务失败", userId, topicId, err)
	}
}
func (u TopicService) UpdateTopicListShowCount(list []entity.Topic, userId int64) {
	err := initUserFlowPool.Submit(func() {
		for _, topic := range list {
			DefaultTopicFlowService.AddUserFlowShowCount(userId, topic.Id)
		}
	})
	if err != nil {
		app.Logger.Error("提交更新topic曝光次数任务失败", userId, err)
	}
}
func (u TopicService) sortTopicListByIds(list []entity.Topic, ids []int64) []entity.Topic {
	topicMap := make(map[int64]entity.Topic)
	for _, topic := range list {
		topicMap[topic.Id] = topic
	}

	newList := make([]entity.Topic, 0)
	for _, id := range ids {
		newList = append(newList, topicMap[id])
	}
	return newList
}
func (u TopicService) GetShareWeappQrCode(userId int, topicId int) ([]byte, string, error) {
	resp, err := wxapp.NewClient(app.Weapp).GetUnlimitedQRCodeResponse(&weapp.UnlimitedQRCode{
		Scene:     fmt.Sprintf("topicId=%d&userId=%d", topicId, userId),
		Page:      "pages/cool-mio/mio-detail/index",
		Width:     100,
		IsHyaline: true,
	})
	if err != nil {
		return nil, "", err
	}
	if resp.ErrCode != 0 {
		return nil, "", errors.New(resp.ErrMsg)
	}
	return resp.Buffer, resp.ContentType, nil
}
func (u TopicService) FindById(topicId int64) entity.Topic {
	return u.r.FindById(topicId)
}
