package service

import (
	contextMix "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

var DefaultUserExtendService = NewUserExtendService(
	repository.NewUserExtendInfoRepository(contextMix.NewMioContext()))

func NewUserExtendService(rUserExtend repository.UserExtentInfoRepository) UserExtendService {
	return UserExtendService{
		rUserExtend: rUserExtend,
	}
}

type UserExtendService struct {
	rUserExtend repository.UserExtentInfoRepository
}

// GetUserExtend 查询用户信息

func (u UserExtendService) GetUserExtend(by repository.GetUserExtendBy) (*entity.UserExtendInfo, bool, error) {
	return u.rUserExtend.GetUserExtend(by)
}
