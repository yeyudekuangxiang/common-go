package business

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/business"
)

var DefaultPointCollectHistoryRepository = PointCollectHistoryRepository{DB: app.DB}

type PointCollectHistoryRepository struct {
	DB *gorm.DB
}

func (repo PointCollectHistoryRepository) Create(history *business.PointCollectHistory) error {
	return repo.DB.Create(history).Error
}
func (repo PointCollectHistoryRepository) Save(history *business.PointCollectHistory) error {
	return repo.DB.Save(history).Error
}
