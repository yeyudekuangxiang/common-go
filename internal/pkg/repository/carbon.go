package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
)

func NewCarbonRepository(ctx *context.MioContext) CarbonRepository {
	return CarbonRepository{ctx: ctx}
}

type CarbonRepository struct {
	ctx *context.MioContext
}

func (repo CarbonRepository) Save(carbon *entity.Carbon) error {
	app.DB.Session(&gorm.Session{})
	return repo.ctx.DB.Save(carbon).Error
}

func (repo CarbonRepository) FindBy(by FindCarbonBy) entity.Carbon {
	carbon := entity.Carbon{}
	db := repo.ctx.DB.Model(carbon)
	if by.OpenId != "" {
		db.Where("openid = ?", by.OpenId)
	}
	if err := db.First(&carbon).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return carbon
}

func (repo CarbonRepository) FindForUpdate(openId string) (entity.Carbon, error) {
	carbon := entity.Carbon{}
	err := repo.ctx.DB.
		Set("gorm:query_option", "for update").
		Where("openid = ?", openId).
		First(&carbon).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return carbon, err
	}
	return carbon, nil
}
func (repo CarbonRepository) Create(carbon *entity.Carbon) error {
	return repo.ctx.DB.Create(carbon).Error
}
