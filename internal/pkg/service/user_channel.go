package service

import (
	"errors"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

var DefaultUserChannelService = NewUserChannelService(repository.DefaultUserChannelRepository)

func NewUserChannelService(channel repository.UserChannelRepository) UserChannelService {
	return UserChannelService{r: channel}
}

type UserChannelService struct {
	r repository.UserChannelRepository
}

func (srv UserChannelService) Create(param *entity.UserChannel) error {
	channel, _ := srv.GetByCid(param.Cid)
	if channel.Cid != 0 {
		return errors.New("渠道已存在，不能重复创建")
	}
	ch := srv.r.FindByCode(repository.FindUserChannelBy{Cid: param.Cid, Code: param.Code})
	if ch.Code != "" {
		return errors.New("code已存在，换一个吧")
	}
	return srv.r.Create(param)
}

func (srv UserChannelService) Save(param *entity.UserChannel) error {
	return srv.r.Save(param)
}

/**修改渠道信息*/
func (srv UserChannelService) UpdateUserChannel(params *entity.UserChannel) error {
	channel := srv.r.FindByCid(repository.FindUserChannelBy{Cid: params.Cid})
	if channel.Cid == 0 {
		return errors.New("渠道不存在，不能修改")
	}
	ch := srv.r.FindByCode(repository.FindUserChannelBy{Cid: params.Cid, Code: params.Code})
	if ch.Code != "" {
		return errors.New("code已存在")
	}
	channel.Code = params.Code
	channel.Name = params.Name
	channel.UpdateTime = model.NewTime()

	err := srv.Save(channel)
	return err
}

/**根据cid获取渠道信息*/
func (srv UserChannelService) GetChannelByCid(cid int64) int64 {
	//如果cid不传那么默认自然流量
	if cid == 0 {
		return 1
	}
	channel, err := srv.GetByCid(cid)
	if err != nil {
		return 1 //不明来源的渠道，算自然流露
	}
	return channel.Cid
}

/**根据cid获取渠道信息*/
func (srv UserChannelService) GetByCid(cid int64) (channel *entity.UserChannel, err error) {
	ch := srv.r.FindByCid(repository.FindUserChannelBy{
		Cid: cid,
	})
	if ch.Cid != 0 {
		return ch, nil
	}
	return ch, errors.New("渠道不存在")
}

/**获取渠道列表*/
func (srv UserChannelService) GetUserChannelPageList(by repository.GetUserChannelPageListBy) ([]entity.UserChannel, int64, error) {
	list, total := srv.r.GetUserChannelPageList(by)
	return list, total, nil
}
