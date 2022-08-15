package business

import (
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity/business"
)

type AreaRepository struct {
	ctx *context.MioContext
}

func NewAreaRepository(ctx *context.MioContext) *AreaRepository {
	return &AreaRepository{ctx: ctx}
}

// List 获取区域列表
func (repo AreaRepository) List(param AreaListPO) ([]business.Area, error) {
	db := repo.ctx.DB.Model(business.Area{})

	if len(param.CityCodes) > 0 {
		db.Where("city_code in (?)", param.CityCodes)
	}

	if len(param.CityIds) > 0 {
		db.Where("city_id in (?)", param.CityIds)
	}

	if param.LikeName != "" {
		db.Where("name like ?", "%"+param.LikeName+"%")
	}
	if param.LikePy != "" {
		db.Where("py like ?", "%"+param.LikePy+"%")
	}
	if param.LikeShortPy != "" {
		db.Where("short_py like ?", "%"+param.LikeShortPy+"%")
	}

	if param.Level != "" {
		db.Where("level = ?", param.Level)
	}

	if param.Search != "" {
		db.Where("name like ? or py like ? or short_py like ?", "%"+param.Search+"%", "%"+param.Search+"%", "%"+param.Search+"%")
	}

	list := make([]business.Area, 0)
	return list, db.Find(&list).Error
}
