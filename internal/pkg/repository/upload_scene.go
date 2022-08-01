package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultUploadSceneRepository = UploadSceneRepository{DB: app.DB}

type UploadSceneRepository struct {
	DB *gorm.DB
}

func (repo UploadSceneRepository) FindScene(by FindSceneBy) (*entity.UploadScene, error) {
	scene := entity.UploadScene{}

	db := repo.DB.Model(scene)
	if by.Scene != "" {
		db.Where("scene = ?", by.Scene)
	}
	return &scene, db.Take(&scene).Error
}
