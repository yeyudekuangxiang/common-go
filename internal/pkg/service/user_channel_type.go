package service

import (
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

var DefaultUserChannelTypeService = NewUserChannelTypeService(repository.DefaultUserChannelTypeRepository)

func NewUserChannelTypeService(channel repository.UserChannelTypeRepository) UserChannelTypeService {
	return UserChannelTypeService{r: channel}
}

type UserChannelTypeService struct {
	r repository.UserChannelTypeRepository
}

func (srv UserChannelTypeService) Create(param *entity.UserChannelType) error {
	return srv.r.Create(param)
}
