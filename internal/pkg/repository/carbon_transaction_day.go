package repository

import (
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository/repotypes"
)

func NewCarbonTransactionDayRepository(ctx *context.MioContext) CarbonTransactionDayRepository {
	return CarbonTransactionDayRepository{
		ctx: ctx,
	}
}

type CarbonTransactionDayRepository struct {
	ctx *context.MioContext
}

func (repo CarbonTransactionDayRepository) Save(day *entity.CarbonTransactionDay) error {
	return repo.ctx.DB.Save(day).Error
}

func (repo CarbonTransactionDayRepository) Create(day *entity.CarbonTransactionDay) error {
	return repo.ctx.DB.Create(day).Error
}

func (repo CarbonTransactionDayRepository) GetList(by repotypes.GetCarbonTransactionDayGetListDO) ([]entity.CarbonTransactionDay, error) {
	list := make([]entity.CarbonTransactionDay, 0)
	db := repo.ctx.DB.Model(entity.CarbonTransactionDay{})
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
