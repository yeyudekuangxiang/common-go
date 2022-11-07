package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
)

func NewUserExtendInfoRepository(ctx *context.MioContext) UserExtentInfoRepository {
	return UserExtentInfoRepository{ctx: ctx}
}

type UserExtentInfoRepository struct {
	ctx *context.MioContext
}

func (repo UserExtentInfoRepository) Save(userExtend *entity.UserExtendInfo) error {
	return repo.ctx.DB.Save(userExtend).Error
}

func (repo UserExtentInfoRepository) GetUserExtend(by GetUserExtendBy) (*entity.UserExtendInfo, bool, error) {
	userExtend := entity.UserExtendInfo{}
	db := repo.ctx.DB.Model(userExtend)
	if by.OpenId != "" {
		db.Where("openid = ?", by.OpenId)
	}
	if by.Uid != 0 {
		db.Where("uid = ?", by.Uid)
	}
	err := db.First(&userExtend).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &userExtend, true, nil
}

func (repo UserExtentInfoRepository) FindForUpdate(openId string) (entity.UserExtendInfo, error) {
	carbon := entity.UserExtendInfo{}
	err := repo.ctx.DB.
		Set("gorm:query_option", "for update").
		Where("openid = ?", openId).
		First(&carbon).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return carbon, err
	}
	return carbon, nil
}

func (repo UserExtentInfoRepository) Create(userExtend *entity.UserExtendInfo) error {
	return repo.ctx.DB.Create(userExtend).Error
}
