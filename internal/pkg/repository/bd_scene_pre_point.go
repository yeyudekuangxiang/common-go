package repository

import (
	"gorm.io/gorm"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
)

var DefaultBdScenePrePointRepository = NewBdScenePrePointModel(mioContext.NewMioContext())

type (
	BdScenePrePointModel interface {
		FindByPlatformUser(platformKey, memberId string) ([]entity.BdScenePrePoint, int64, error)
		FindByOpenId(openId, platformKey string) ([]entity.BdScenePrePoint, int64, error)
		FindAllByOpenId(openId string) ([]entity.BdScenePrePoint, int64, error)
		FindAllByPlatformKey(platformKey string) ([]entity.BdScenePrePoint, int64, error)
		FindBy(by GetScenePrePoint) ([]entity.BdScenePrePoint, int64, error)
		FindOne(by GetScenePrePoint) (entity.BdScenePrePoint, bool, error)
		Create(data *entity.BdScenePrePoint) error
		Save(data *entity.BdScenePrePoint) error
		Updates(cond GetScenePrePoint, up map[string]interface{}) error
	}
	defaultBdScenePrePointModel struct {
		ctx *mioContext.MioContext
	}
)

func NewBdScenePrePointModel(ctx *mioContext.MioContext) BdScenePrePointModel {
	return &defaultBdScenePrePointModel{
		ctx: ctx,
	}
}

func (m defaultBdScenePrePointModel) FindByPlatformUser(platformKey, memberId string) ([]entity.BdScenePrePoint, int64, error) {
	var item []entity.BdScenePrePoint
	var total int64
	err := m.ctx.DB.
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

func (m defaultBdScenePrePointModel) FindByOpenId(openId, platformKey string) ([]entity.BdScenePrePoint, int64, error) {
	var item []entity.BdScenePrePoint
	var total int64
	err := m.ctx.DB.
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

func (m defaultBdScenePrePointModel) FindAllByOpenId(openId string) ([]entity.BdScenePrePoint, int64, error) {
	var items []entity.BdScenePrePoint
	var total int64
	err := m.ctx.DB.
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

func (m defaultBdScenePrePointModel) FindAllByPlatformKey(platformKey string) ([]entity.BdScenePrePoint, int64, error) {
	var items []entity.BdScenePrePoint
	var total int64
	err := m.ctx.DB.Where("platform_key = ?", platformKey).Count(&total).Find(&items).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
		return []entity.BdScenePrePoint{}, 0, err
	}
	return items, total, nil
}

func (m defaultBdScenePrePointModel) FindBy(by GetScenePrePoint) ([]entity.BdScenePrePoint, int64, error) {
	var items []entity.BdScenePrePoint
	query := m.ctx.DB.Model(&entity.BdScenePrePoint{}).Where("platform_key = ?", by.PlatformKey)
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

func (m defaultBdScenePrePointModel) FindOne(by GetScenePrePoint) (entity.BdScenePrePoint, bool, error) {
	var item entity.BdScenePrePoint

	query := m.ctx.DB.Where("id = ?", by.Id)
	if by.OpenId != "" {
		query = query.Where("open_id = ?", by.OpenId)
	}
	if by.PlatformUserId != "" {
		query = query.Where("platform_user_id = ?", by.PlatformUserId)
	}
	if by.PlatformKey != "" {
		query = query.Where("platform_key = ?", by.PlatformKey)
	}
	if by.Status != 0 {
		query = query.Where("status = ?", by.Status)
	}
	if by.TradeNo != "" {
		query = query.Where("tradeno = ?", by.TradeNo)
	}

	err := query.First(&item).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
		return entity.BdScenePrePoint{}, false, nil
	}
	return item, true, nil
}

func (m defaultBdScenePrePointModel) Create(data *entity.BdScenePrePoint) error {
	return m.ctx.DB.Model(&entity.BdScenePrePoint{}).Create(data).Error
}

func (m defaultBdScenePrePointModel) Save(data *entity.BdScenePrePoint) error {
	return m.ctx.DB.Save(data).Error
}

func (m defaultBdScenePrePointModel) Updates(cond GetScenePrePoint, up map[string]interface{}) error {
	query := m.ctx.DB.Model(&entity.BdScenePrePoint{})
	if cond.Id != 0 {
		query.Where("id = ?", cond.Id)
	}

	if len(cond.Ids) > 0 {
		query.Where("id in (?)", cond.Ids)
	}

	if cond.OpenId != "" {
		query.Where("open_id = ?", cond.OpenId)
	}

	if cond.PlatformUserId != "" {
		query.Where("platform_user_id = ?", cond.PlatformUserId)
	}

	if cond.PlatformKey != "" {
		query.Where("platform_key = ?", cond.PlatformKey)
	}

	if cond.Status != 0 {
		query.Where("status = ?", cond.Status)
	}

	err := query.Updates(up).Error

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
		return err
	}
	return nil
}
