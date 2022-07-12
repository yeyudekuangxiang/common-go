package service

import (
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

var DefaultUserRiskLogService = UserRiskLogService{}

type UserRiskLogService struct {
}

func (srv UserRiskLogService) Create(param *entity.UserRiskLog) error {
	return repository.DefaultUserRiskLogRepository.Create(param)
}
