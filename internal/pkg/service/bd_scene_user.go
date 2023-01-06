package service

import (
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/pkg/errno"
)

var DefaultBdSceneUserService = BdSceneUserService{}

type BdSceneUserService struct {
}

func (srv BdSceneUserService) FindOne(params repository.GetSceneUserOne) *entity.BdSceneUser {
	return repository.DefaultBdSceneUserRepository.FindOne(params)
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
	sceneUser := srv.FindOne(repository.GetSceneUserOne{
		PlatformKey:    scene.Ch,
		PlatformUserId: memberId,
	})

	if sceneUser.ID != 0 {
		return sceneUser, errno.ErrExisting
	}

	sceneUser.PlatformKey = scene.Ch
	sceneUser.PlatformUserId = memberId
	sceneUser.Phone = user.PhoneNumber
	sceneUser.OpenId = user.OpenId
	sceneUser.UnionId = user.UnionId
	err := srv.Create(sceneUser)
	if err != nil {
		return nil, err
	}
	return sceneUser, nil
}
