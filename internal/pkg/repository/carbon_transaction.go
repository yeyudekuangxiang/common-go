package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository/repotypes"
)

func NewCarbonTransactionRepository(ctx *context.MioContext) CarbonTransactionRepository {
	return CarbonTransactionRepository{ctx: ctx}
}

type CarbonTransactionRepository struct {
	ctx *context.MioContext
}

func (repo CarbonTransactionRepository) Save(transaction *entity.CarbonTransaction) error {
	return repo.ctx.DB.Save(transaction).Error
}

func (repo CarbonTransactionRepository) Create(transaction *entity.CarbonTransaction) error {
	return repo.ctx.DB.Create(transaction).Error
}

func (repo CarbonTransactionRepository) GetListBy(by repotypes.GetCarbonTransactionListByDO) []repotypes.GetCarbonTransactionListBy {
	list := make([]repotypes.GetCarbonTransactionListBy, 0)
	db := repo.ctx.DB.Model(entity.CarbonTransaction{})
	db.Select("type", "sum(value)", "user_id")
	if by.OpenId != "" {
		db.Where("openid = ?", by.OpenId)
	}
	if by.StartTime == "" {
		db.Where("created_at >= ?", by.StartTime)
	}
	if by.EndTime != "" {
		db.Where("created_at <= ?", by.EndTime)
	}
	if by.Type != "" {
		db.Where("type = ?", by.Type)
	}
	db.Group("type,user_id")
	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return list
}

func (repo CarbonTransactionRepository) GetListByDay(by repotypes.GetCarbonTransactionListByDO) []repotypes.GetCarbonTransactionListBy {
	list := make([]repotypes.GetCarbonTransactionListBy, 0)
	db := repo.ctx.DB.Model(entity.CarbonTransaction{})
	db.Select("sum(value)", "user_id", "openid")
	if by.StartTime == "" {
		db.Where("created_at >= ?", by.StartTime)
	}
	if by.EndTime != "" {
		db.Where("created_at <= ?", by.EndTime)
	}
	db.Group("user_id")
	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return list
}

func (repo CarbonTransactionRepository) FindBy(by FindPointTransactionBy) entity.CarbonTransaction {
	pt := entity.CarbonTransaction{}
	db := repo.ctx.DB.Model(pt)
	if by.TransactionId != "" {
		db.Where("transaction_id", by.TransactionId)
	}
	err := db.First(&pt).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return pt
}
func (repo CarbonTransactionRepository) GetPageListBy(by GetPointTransactionPageListBy) ([]entity.CarbonTransaction, int64) {
	list := make([]entity.CarbonTransaction, 0)

	db := repo.ctx.DB.Model(entity.CarbonTransaction{})
	if len(by.OpenIds) > 0 {
		db.Where("openid in (?)", by.OpenIds)
	}

	if !by.StartTime.IsZero() {
		db.Where("create_time >= ?", by.StartTime.Time)
	}
	if !by.EndTime.IsZero() {
		db.Where("create_time <= ?", by.EndTime.Time)
	}

	if by.Type != "" {
		db.Where("type = ?", by.Type)
	}
	if len(by.Types) > 0 {
		db.Where("type in (?)", by.Types)
	}
	if by.AdminId != 0 {
		db.Where("admin_id = ?", by.AdminId)
	}

	for _, orderBy := range by.OrderBy {
		switch orderBy {
		case entity.OrderByPointTranCTDESC:
			db.Order("create_time desc")
		}
	}

	var total int64
	db.Count(&total).Limit(by.Limit).Offset(by.Offset)

	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}

	return list, total
}
