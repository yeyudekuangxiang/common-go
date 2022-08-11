package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"strings"
)

var DefaultPointCollectHistoryRepository = PointCollectHistoryRepository{DB: app.DB}

type PointCollectHistoryRepository struct {
	DB *gorm.DB
}

func (repo PointCollectHistoryRepository) Create(history *entity.PointCollectHistory) error {
	return repo.DB.Create(history).Error
}

func (repo PointCollectHistoryRepository) Count(history *entity.PointCollectHistory) (count int64, err error) {
	model := repo.DB.Model(&entity.PointCollectHistory{})
	if history.OpenId != "" {
		model.Where("openid = ?", history.OpenId)
	}
	if history.Type != "" {
		model.Where("type = ?", strings.ToUpper(history.Type))
	}
	if !history.Date.IsZero() {
		model.Where("date = ?", history.Date)
	}
	if !history.Time.IsZero() {
		model.Where("time = ?", history.Time)
	}
	if err = model.Count(&count).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return
}

func (repo PointCollectHistoryRepository) CreateLog(history *entity.PointCollectLog) error {
	return repo.DB.Create(history).Error
}
