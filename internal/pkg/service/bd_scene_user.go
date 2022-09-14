package service

import (
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

var DefaultBdSceneUserService = BdSceneUserService{}

type BdSceneUserService struct {
}

func (srv BdSceneUserService) FindByCh(ch string) *entity.BdSceneUser {
	item := repository.DefaultBdSceneUserRepository.FindByCh(ch)
	return &item
}

func (srv BdSceneUserService) Create(data *entity.BdSceneUser) error {
	return repository.DefaultBdSceneUserRepository.Create(data)
}
