package service

import (
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service/platform/jhx"
)

var DefaultBdSceneUserService = BdSceneUserService{}

type BdSceneUserService struct {
}

func (srv BdSceneUserService) FindByCh(platformKey string) *entity.BdSceneUser {
	item := repository.DefaultBdSceneUserRepository.FindByCh(platformKey)
	return &item
}

func (srv BdSceneUserService) FindPlatformUser(openId string, platformKey string) *entity.BdSceneUser {
	item := repository.DefaultBdSceneUserRepository.FindPlatformUser(openId, platformKey)
	return &item
}

func (srv BdSceneUserService) FindPlatformUserByPlatformUserId(memberId string, platformKey string) *entity.BdSceneUser {
	item := repository.DefaultBdSceneUserRepository.FindPlatformUserByPlatformUserId(memberId, platformKey)
	return &item
}

func (srv BdSceneUserService) FindPlatformUserByOpenId(openId string) *entity.BdSceneUser {
	item := repository.DefaultBdSceneUserRepository.FindPlatformUserByOpenId(openId)
	return &item
}

func (srv BdSceneUserService) Create(data *entity.BdSceneUser) error {
	return repository.DefaultBdSceneUserRepository.Create(data)
}

func (srv BdSceneUserService) Bind(user entity.User, scene entity.BdScene, memberId string) (*entity.BdSceneUser, error) {
	sceneUser := srv.FindPlatformUser(user.OpenId, scene.Ch)
	if sceneUser.ID != 0 {
		return sceneUser, nil
	}
	sceneUser.PlatformKey = scene.Ch
	sceneUser.PlatformUserId = memberId
	sceneUser.Phone = user.PhoneNumber
	sceneUser.OpenId = user.OpenId
	sceneUser.UnionId = user.UnionId
	err := srv.Create(sceneUser)
	if err != nil {
		app.Logger.Errorf("create db_scene_user error:%s", err.Error())
		return nil, err
	}
	//绑定回调
	if scene.Ch == "jinhuaxing" {
		err = jhx.NewJhxService(context.NewMioContext()).BindSuccess(sceneUser.Phone, "1")
		if err != nil {
			app.Logger.Errorf("callback %s error:%s", scene.Ch, err.Error())
			return nil, err
		}
	}
	return sceneUser, nil
}
