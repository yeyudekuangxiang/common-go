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

func (repo BdScenePrePointRepository) FindByPlatformUser(memberId string, platformKey string) ([]entity.BdScenePrePoint, int64, error) {
	var item []entity.BdScenePrePoint
	var total int64
	err := repo.DB.
		Where("platform_key = ?", platformKey).
		Where("platform_user_id = ?", memberId).
		Count(&total).
		Find(&item).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
		return []entity.BdScenePrePoint{}, 0, err
	}
	return item, total, nil
}

func (repo BdScenePrePointRepository) FindByOpenId(openId, platformKey string) ([]entity.BdScenePrePoint, int64, error) {
	var item []entity.BdScenePrePoint
	var total int64
	err := repo.DB.
		Where("platform_key = ?", platformKey).
		Where("open_id = ?", openId).
		Count(&total).
		Find(&item).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
		return []entity.BdScenePrePoint{}, 0, err
	}
	return item, total, nil
}

func (repo BdScenePrePointRepository) FindAllByOpenId(openId string) ([]entity.BdScenePrePoint, int64, error) {
	var items []entity.BdScenePrePoint
	var total int64
	err := repo.DB.
		Where("open_id = ?", openId).
		Count(&total).
		Find(&items).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
		return []entity.BdScenePrePoint{}, 0, err
	}
	return items, total, nil
}

func (repo BdScenePrePointRepository) FindAllByPlatformKey(platformKey string) ([]entity.BdScenePrePoint, int64, error) {
	var items []entity.BdScenePrePoint
	var total int64
	err := repo.DB.Where("platform_key = ?", platformKey).Count(&total).Find(&items).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
		return []entity.BdScenePrePoint{}, 0, err
	}
	return items, total, nil
}

func (repo BdScenePrePointRepository) FindBy(by GetScenePrePoint) ([]entity.BdScenePrePoint, int64, error) {
	var items []entity.BdScenePrePoint
	query := repo.DB.Model(&entity.BdScenePrePoint{}).Where("platform_key = ?", by.PlatformKey)
	if by.OpenId != "" {
		query.Where("open_id = ?", by.OpenId)
	}
	if by.PlatformUserId != "" {
		query.Where("platform_user_id = ?", by.PlatformUserId)
	}
	if !by.StartTime.IsZero() {
		query.Where("created_at > ?", by.StartTime)
	}
	if !by.EndTime.IsZero() {
		query.Where("created_at < ?", by.EndTime)
	}
	if by.Status != 0 {
		query.Where("status = ?", by.Status)
	}
	var total int64
	err := query.Count(&total).Order("created_at asc").Find(&items).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
		return []entity.BdScenePrePoint{}, 0, err
	}
	return items, total, nil
}

func (repo BdScenePrePointRepository) FindOne(by GetScenePrePoint) (entity.BdScenePrePoint, error) {
	var item entity.BdScenePrePoint

	query := repo.DB.Where("id = ?", by.Id)
	if by.OpenId != "" {
		query.Where("open_id = ?", by.OpenId)
	}
	if by.PlatformUserId != "" {
		query.Where("platform_user_id = ?", by.PlatformUserId)
	}
	if by.PlatformKey != "" {
		query.Where("platform_key = ?", by.PlatformKey)
	}
	if by.Status != 0 {
		query.Where("status = ?", by.Status)
	}

	err := query.First(&item).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
		return entity.BdScenePrePoint{}, err
	}
	return item, nil
}

func (repo BdScenePrePointRepository) Create(data *entity.BdScenePrePoint) error {
	return repo.DB.Model(&entity.BdScenePrePoint{}).Create(data).Error
}

func (repo BdScenePrePointRepository) Save(data *entity.BdScenePrePoint) error {
	return repo.DB.Save(data).Error
}
