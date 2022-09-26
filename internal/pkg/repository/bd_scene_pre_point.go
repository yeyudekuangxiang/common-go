package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultBdScenePrePointRepository = BdScenePrePointRepository{DB: app.DB}

type BdScenePrePointRepository struct {
	DB *gorm.DB
}

func (repo BdScenePrePointRepository) FindByPlatformUser(memberId string, platformKey string) entity.BdScenePrePoint {
	item := entity.BdScenePrePoint{}
	err := repo.DB.
		Where("platform_key = ?", platformKey).
		Where("platform_user_id = ?", memberId).
		First(&item).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return item
}

func (repo BdScenePrePointRepository) FindByOpenId(openId, platformKey string) entity.BdScenePrePoint {
	item := entity.BdScenePrePoint{}
	err := repo.DB.
		Where("platform_key = ?", platformKey).
		Where("open_id = ?", openId).
		First(&item).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return item
}

func (repo BdScenePrePointRepository) FindAllByOpenId(openId string) []entity.BdScenePrePoint {
	var items []entity.BdScenePrePoint
	err := repo.DB.
		Where("open_id = ?", openId).
		Find(&items).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return items
}

func (repo BdScenePrePointRepository) FindAllByPlatformKey(platformKey string) []entity.BdScenePrePoint {
	var items []entity.BdScenePrePoint
	err := repo.DB.Where("platform_key = ?", platformKey).Find(&items).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return items
}

func (repo BdScenePrePointRepository) FindBy(by GetScenePrePoint) []entity.BdScenePrePoint {
	var items []entity.BdScenePrePoint
	query := repo.DB.Where("platform_key = ?", by.PlatformKey)
	if by.OpenId != "" {
		query.Where("open_id = ?", by.OpenId)
	}
	if by.PlatformUserId != "" {
		query.Where("platform_user_id = ?", by.PlatformUserId)
	}
	if by.Status != 0 {
		query.Where("status = ?", by.Status)
	}
	if !by.StartTime.IsZero() {
		query.Where("to_char(create_time,'YYYY-MM-DD') < ?", by.StartTime.Format("2006-01-02"))
	}
	if !by.EndTime.IsZero() {
		query.Where("to_char(create_time,'YYYY-MM-DD') > ?", by.EndTime.Format("2006-01-02"))
	}
	err := query.Find(&items).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return items
}

func (repo BdScenePrePointRepository) Create(data entity.BdScenePrePoint) error {
	return repo.DB.Model(&entity.BdScenePrePoint{}).Create(&data).Error
}
