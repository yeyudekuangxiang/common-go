package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultBdSceneUserRepository = BdSceneUserRepository{DB: app.DB}

type BdSceneUserRepository struct {
	DB *gorm.DB
}

func (repo BdSceneUserRepository) FindByCh(platformKey string) entity.BdSceneUser {
	item := entity.BdSceneUser{}
	err := repo.DB.Where("platform_key = ?", platformKey).First(&item).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return item
}

func (repo BdSceneUserRepository) FindPlatformUser(openId string, platformKey string) entity.BdSceneUser {
	item := entity.BdSceneUser{}
	err := repo.DB.
		Where("open_id = ?", openId).
		Where("platform_key = ?", platformKey).
		First(&item).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return item
}

func (repo BdSceneUserRepository) FindPlatformUserByOpenId(openId string) entity.BdSceneUser {
	item := entity.BdSceneUser{}
	err := repo.DB.
		Where("open_id = ?", openId).
		First(&item).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return item
}

func (repo BdSceneUserRepository) Create(data *entity.BdSceneUser) error {
	return repo.DB.Model(&entity.BdSceneUser{}).Create(data).Error
}
