package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultUserRiskLogRepository = UserRiskLogRepository{DB: app.DB}

type UserRiskLogRepository struct {
	DB *gorm.DB
}

func (repo UserRiskLogRepository) Create(param *entity.UserRiskLog) error {
	return repo.DB.Create(param).Error
}
