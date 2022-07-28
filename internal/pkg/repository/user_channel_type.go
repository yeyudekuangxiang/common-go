package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultUserChannelTypeRepository = UserChannelTypeRepository{DB: app.DB}

type UserChannelTypeRepository struct {
	DB *gorm.DB
}

func (repo UserChannelTypeRepository) Save(channel *entity.UserChannelType) error {
	return repo.DB.Save(channel).Error
}

func (repo UserChannelTypeRepository) Create(channel *entity.UserChannelType) error {
	return repo.DB.Create(channel).Error
}
