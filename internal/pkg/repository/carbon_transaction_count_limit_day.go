package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository/repotypes"
)

func NewCarbonTransactionCountLimitRepository(ctx *context.MioContext) CarbonTransactionCountLimitRepository {
	return CarbonTransactionCountLimitRepository{ctx: ctx}
}

type CarbonTransactionCountLimitRepository struct {
	ctx *context.MioContext
}

func (repo CarbonTransactionCountLimitRepository) Save(transaction *entity.CarbonTransactionCountLimit) error {
	return repo.ctx.DB.Save(transaction).Error
}

func (repo CarbonTransactionCountLimitRepository) Create(transaction *entity.CarbonTransactionCountLimit) error {
	return repo.ctx.DB.Create(transaction).Error
}

func (repo CarbonTransactionCountLimitRepository) GetListBy(by repotypes.GetCarbonTransactionCountLimitDO) []entity.CarbonTransactionCountLimit {
	list := make([]entity.CarbonTransactionCountLimit, 0)

	db := repo.ctx.DB.Model(entity.CarbonTransactionCountLimit{})
	if by.OpenId != "" {
		db.Where("openid = ?", by.OpenId)
	}

	if !by.StartTime.IsZero() {
		db.Where("create_time >= ?", by.StartTime)
	}
	if !by.EndTime.IsZero() {
		db.Where("create_time <= ?", by.EndTime)
	}

	if by.Type != "" {
		db.Where("type = ?", by.Type)
	}

	for _, orderBy := range by.OrderBy {
		switch orderBy {
		case entity.OrderByPointTranCTDESC:
			db.Order("create_time desc")
		}
	}

	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}

	return list
}

func (p CarbonTransactionCountLimitRepository) FindBy(by repotypes.FindCarbonTransactionCountLimitFindByDO) entity.CarbonTransactionCountLimit {
	limit := entity.CarbonTransactionCountLimit{}
	db := p.ctx.DB.Model(entity.CarbonTransactionCountLimit{})
	if by.OpenId != "" {
		db.Where("openid = ?", by.OpenId)
	}
	if by.Type != "" {
		db.Where("type = ?", by.Type)
	}
	if by.VDate != "" {
		db.Where("v_date = ?", by.VDate)
	}
	if err := db.First(&limit).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return limit
}
