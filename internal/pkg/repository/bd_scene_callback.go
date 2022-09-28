package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultBdSceneCallbackRepository = BdSceneCallbackRepository{DB: app.DB}

type BdSceneCallbackRepository struct {
	DB *gorm.DB
}

func (repo BdSceneCallbackRepository) FindOne(by GetSceneCallback) entity.BdSceneCallback {
	item := entity.BdSceneCallback{}
	query := repo.DB.Where("platform_key = ?", by.PlatformKey)
	if by.OpenId != "" {
		query.Where("open_id = ?", by.OpenId)
	}
	if by.SourceKey != "" {
		query.Where("source_key = ?", by.SourceKey)
	}
	if by.PlatformUserId != "" {
		query.Where("platform_user_id = ?", by.PlatformUserId)
	}
	err := query.First(&item).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return item
}

func (repo BdSceneCallbackRepository) FindAll(by GetSceneCallback) ([]entity.BdSceneCallback, int64, error) {
	var items []entity.BdSceneCallback
	var total int64
	query := repo.DB.Where("platform_key = ?", by.PlatformKey)
	if by.OpenId != "" {
		query.Where("open_id = ?", by.OpenId)
	}
	if by.SourceKey != "" {
		query.Where("source_key = ?", by.SourceKey)
	}
	if by.PlatformUserId != "" {
		query.Where("platform_user_id = ?", by.PlatformUserId)
	}
	if by.StartTime != "" {
		query.Where("to_char(created_at, 'YYYY-MM-DD HH:MI:SS') > ?", by.StartTime)
	}
	if by.EndTime != "" {
		query.Where("to_char(created_at, 'YYYY-MM-DD HH:MI:SS') < ?", by.EndTime)
	}
	err := query.Count(&total).Find(&items).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return items, 0, err
		}
		panic(err)
	}
	return items, total, nil
}

func (repo BdSceneCallbackRepository) Save(callback entity.BdSceneCallback) error {
	return repo.DB.Model(&entity.BdSceneCallback{}).Save(&callback).Error
}
