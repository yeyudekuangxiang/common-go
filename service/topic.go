package service

import (
	"fmt"
	"github.com/medivhzhan/weapp/v3"
	"github.com/mlogclub/simple"
	"github.com/pkg/errors"
	"mio/core/app"
	"mio/internal/wxapp"
	"mio/model"
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

func (u TopicService) List(cnq *simple.SqlCnd) (list []model.Topic) {
	return u.r.List(cnq)
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
