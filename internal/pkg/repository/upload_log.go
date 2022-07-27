package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultUploadLogRepository = UploadLogRepository{DB: app.DB}

type UploadLogRepository struct {
	DB *gorm.DB
}

func (repo UploadLogRepository) FindLog(by FindLogBy) (*entity.UploadLog, error) {
	Log := entity.UploadLog{}

	db := repo.DB.Model(Log)
	if by.SceneId != 0 {
		db.Where("scene_id = ?", by.SceneId)
	}
	if by.LogId != "" {
		db.Where("log_id = ?", by.LogId)
	}
	return &Log, db.Take(&Log).Error
}
func (repo UploadLogRepository) Create(log *entity.UploadLog) error {
	return repo.DB.Create(log).Error
}
func (repo UploadLogRepository) Save(log *entity.UploadLog) error {
	return repo.DB.Save(log).Error
}
