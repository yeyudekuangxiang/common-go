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

func (repo BdSceneUserRepository) FindPlatformUserByPlatformUserId(memberId string, platformKey string) entity.BdSceneUser {
	item := entity.BdSceneUser{}
	err := repo.DB.
		Where("platform_key = ?", platformKey).
		Where("platform_user_id = ?", memberId).
		First(&item).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return item
}

func (repo BdSceneUserRepository) FindPlatformUser(openId, platformKey string) entity.BdSceneUser {
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

func (repo BdSceneUserRepository) FindOne(params GetSceneUserOne) *entity.BdSceneUser {
	one := entity.BdSceneUser{}
	query := repo.DB.Model(&one)

	if params.Id != 0 {
		query.Where("id = ?", params.Id)
	}

	if params.PlatformKey != "" {
		query.Where("platform_key = ?", params.PlatformKey)
	}

	if params.PlatformUserId != "" {
		query.Where("platform_user_id = ?", params.PlatformUserId)
	}

	if params.OpenId != "" {
		query.Where("open_id = ?", params.OpenId)
	}

	if params.Phone != "" {
		query.Where("phone = ?", params.Phone)
	}

	if params.UnionId != "" {
		query.Where("union_id = ?", params.UnionId)
	}

	err := query.First(&one).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return &one
}
