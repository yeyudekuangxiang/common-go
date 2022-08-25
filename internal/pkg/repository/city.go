package repository

import (
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository/repotypes"
)

func NewCityRepository(ctx *context.MioContext) CityRepository {
	return CityRepository{
		ctx: ctx,
	}
}

type CityRepository struct {
	ctx *context.MioContext
}

func (repo CityRepository) Save(day *entity.City) error {
	return repo.ctx.DB.Save(day).Error
}

func (repo CityRepository) Create(day *entity.City) error {
	return repo.ctx.DB.Create(day).Error
}

func (repo CityRepository) GetList(by repotypes.GetCarbonTransactionDayGetListDO) ([]entity.City, error) {
	list := make([]entity.City, 0)
	db := repo.ctx.DB.Model(entity.City{})
	if by.StartTime != "" {
		db.Where("v_date >= ?", by.StartTime)
	}
	if by.EndTime != "" {
		db.Where("v_date <= ?", by.EndTime)
	}
	if by.UserId != 0 {
		db.Where("user_id = ?", by.UserId)
	}
	for _, orderBy := range by.OrderBy {
		switch orderBy {
		case entity.OrderByCarbonTranDayVDate:
			db.Order("v_date desc")
		}
	}
	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return list, nil
}
