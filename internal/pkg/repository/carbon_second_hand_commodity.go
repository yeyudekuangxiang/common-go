package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
)

type (
	CarbonSecondHandCommodityModel interface {
		FindOne(id int64) (*entity.CarbonSecondHandCommodity, error)
		FindListByIds(ids []int64) []*entity.CarbonSecondHandCommodity
	}

	defaultCarbonSecondHandCommodityModel struct {
		ctx *mioContext.MioContext
	}
)

func (m *defaultCarbonSecondHandCommodityModel) FindListByIds(ids []int64) []*entity.CarbonSecondHandCommodity {
	commentList := make([]*entity.CarbonSecondHandCommodity, len(ids))
	err := app.DB.Model(&entity.CarbonSecondHandCommodity{}).
		Where("id in (?)", ids).
		//Where("state = ?", 0).
		Find(&commentList).Error
	if err != nil {
		return []*entity.CarbonSecondHandCommodity{}
	}
	return commentList
}

func NewCarbonSecondHandCommodityModel(ctx *mioContext.MioContext) CarbonSecondHandCommodityModel {
	return &defaultCarbonSecondHandCommodityModel{
		ctx: ctx,
	}
}

func (m *defaultCarbonSecondHandCommodityModel) FindOne(id int64) (*entity.CarbonSecondHandCommodity, error) {
	var resp entity.CarbonSecondHandCommodity
	err := m.ctx.DB.First(&resp, id).Error
	switch err {
	case nil:
		return &resp, nil
	case gorm.ErrRecordNotFound:
		return nil, entity.ErrNotFount
	default:
		return nil, err
	}
}
