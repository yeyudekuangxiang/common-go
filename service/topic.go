package service

import (
	"fmt"
	"github.com/medivhzhan/weapp/v3"
	"github.com/pkg/errors"
	"mio/core/app"
	"mio/internal/wxapp"
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

func (u TopicService) GetTopicDetailPageList(param repository.GetTopicPageListBy) ([]TopicDetail, int64, error) {
	list, total := u.r.GetTopicPageList(param)

	//查询点赞信息
	topicIds := make([]int64, 0)
	for _, topic := range list {
		topicIds = append(topicIds, topic.Id)
	}
	topicLikeMap := make(map[int64]bool)
	if param.UserId > 0 {
		likeList := repository.TopicLikeRepository{DB: app.DB}.GetListBy(repository.GetTopicLikeListBy{
			TopicIds: topicIds,
			UserId:   param.UserId,
		})
		for _, like := range likeList {
			topicLikeMap[int64(like.TopicId)] = like.Status == 1
		}
	}

	//整理数据
	detailList := make([]TopicDetail, 0)
	for _, topic := range list {
		detailList = append(detailList, TopicDetail{
			Topic:  topic,
			IsLike: topicLikeMap[topic.Id],
		})
	}

	return detailList, total, nil
}
func (u TopicService) GetShareWeappQrCode(openid string, topicId int) ([]byte, string, error) {
	userRes := repository.NewUserRepository()
	user := userRes.GetShortUserBy(repository.GetUserBy{
		OpenId: openid,
	})
	resp, err := wxapp.NewClient(app.Weapp).GetUnlimitedQRCode(&weapp.UnlimitedQRCode{
		Scene:     fmt.Sprintf("topicId=%d&userId=%d", topicId, user.ID),
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
