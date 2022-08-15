package business

import (
	"fmt"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity/business"
)

type AreaRepository struct {
	ctx *context.MioContext
}

func NewAreaRepository(ctx *context.MioContext) *AreaRepository {
	return &AreaRepository{ctx: ctx}
}

// CityList 获取城市列表
func (repo AreaRepository) CityList(param CityLisPO) ([]business.Area, error) {
	db := repo.ctx.DB.Table("business_area as province")

	db.Joins("business_area as city on province.city_code = city.pid_code").Where("province.pid_code = '86'")

	if param.LikeName != "" {
		db.Where("city.name like ?", fmt.Sprintf(`"%s"`, param.LikeName))
	}
	if param.LikePy != "" {
		db.Where("city.py like ?", fmt.Sprintf(`"%s"`, param.LikePy))
	}
	if param.LikeShortPy != "" {
		db.Where("city.short_py like ?", fmt.Sprintf(`"%s"`, param.LikeShortPy))
	}

	list := make([]business.Area, 0)
	return list, db.Select("city.*").Find(&list).Error
}

// List 获取区域列表
func (repo AreaRepository) List(param AreaListPO) ([]business.Area, error) {
	db := repo.ctx.DB.Model(business.Area{})

	if len(param.CityCodes) > 0 {
		db.Where("city_code in (?)", param.CityCodes)
	}

	list := make([]business.Area, 0)
	return list, db.Find(&list).Error
}
