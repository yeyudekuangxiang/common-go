package service

import (
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
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

func (srv BdSceneUserService) Create(data *entity.BdSceneUser) error {
	return repository.DefaultBdSceneUserRepository.Create(data)
}
