package repository

import (
	"gorm.io/gorm"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
)

var DefaultBdSceneUserRepository = NewBdSceneUserModel(mioContext.NewMioContext())

type (
	BdSceneUserModel interface {
		FindOne(params GetSceneUserOne) *entity.BdSceneUser
		CheckBind(params GetSceneUserOne) *entity.BdSceneUser
		FindByCh(platformKey string) entity.BdSceneUser
		FindPlatformUserByPlatformUserId(memberId string, platformKey string) entity.BdSceneUser
		FindPlatformUser(openId, platformKey string) entity.BdSceneUser
		FindPlatformUserByOpenId(openId string) entity.BdSceneUser
		Create(data *entity.BdSceneUser) error
	}

	defaultBdSceneUserModel struct {
		ctx *mioContext.MioContext
	}
)

func (m defaultBdSceneUserModel) CheckBind(params GetSceneUserOne) *entity.BdSceneUser {
	one := entity.BdSceneUser{}
	query := m.ctx.DB.Model(&one)

	if params.Id != 0 {
		query.Where("id = ?", params.Id)
	}
	if params.PlatformKey == "" || params.PlatformUserId == "" || params.OpenId == "" {
		return &one
	}

	err := query.Where("platform_key = ?", params.PlatformKey).
		Where("platform_user_id = ? or open_id = ?", params.PlatformUserId, params.OpenId).
		First(&one).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return &one
}

func NewBdSceneUserModel(ctx *mioContext.MioContext) BdSceneUserModel {
	return &defaultBdSceneUserModel{
		ctx: ctx,
	}
}

func (m defaultBdSceneUserModel) FindByCh(platformKey string) entity.BdSceneUser {
	item := entity.BdSceneUser{}
	err := m.ctx.DB.Where("platform_key = ?", platformKey).First(&item).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return item
}

func (m defaultBdSceneUserModel) FindPlatformUserByPlatformUserId(memberId string, platformKey string) entity.BdSceneUser {
	item := entity.BdSceneUser{}
	err := m.ctx.DB.
		Where("platform_key = ?", platformKey).
		Where("platform_user_id = ?", memberId).
		First(&item).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return item
}

func (m defaultBdSceneUserModel) FindPlatformUser(openId, platformKey string) entity.BdSceneUser {
	item := entity.BdSceneUser{}
	err := m.ctx.DB.
		Where("open_id = ?", openId).
		Where("platform_key = ?", platformKey).
		First(&item).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return item
}

func (m defaultBdSceneUserModel) FindPlatformUserByOpenId(openId string) entity.BdSceneUser {
	item := entity.BdSceneUser{}
	err := m.ctx.DB.
		Where("open_id = ?", openId).
		First(&item).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return item
}

func (m defaultBdSceneUserModel) Create(data *entity.BdSceneUser) error {
	return m.ctx.DB.Model(&entity.BdSceneUser{}).Create(data).Error
}

func (m defaultBdSceneUserModel) FindOne(params GetSceneUserOne) *entity.BdSceneUser {
	one := entity.BdSceneUser{}
	query := m.ctx.DB.Model(&one)

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
