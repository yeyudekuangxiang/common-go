package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultPointCollectHistoryRepository = PointCollectHistoryRepository{DB: app.DB}

type PointCollectHistoryRepository struct {
	DB *gorm.DB
}

func (repo PointCollectHistoryRepository) Create(history *entity.PointCollectHistory) error {
	return repo.DB.Create(history).Error
}
