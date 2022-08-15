package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
)

func NewPointTransactionRepository(ctx *context.MioContext) *PointTransactionRepository {
	return &PointTransactionRepository{
		ctx: ctx,
	}
}

type PointTransactionRepository struct {
	ctx *context.MioContext
}

func (repo PointTransactionRepository) Save(transaction *entity.PointTransaction) error {
	return repo.ctx.DB.Save(transaction).Error
}

func (repo PointTransactionRepository) GetListBy(by GetPointTransactionListBy) []entity.PointTransaction {
	list := make([]entity.PointTransaction, 0)

	db := repo.ctx.DB.Model(entity.PointTransaction{})
	if by.OpenId != "" {
		db.Where("openid = ?", by.OpenId)
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
func (repo PointTransactionRepository) FindBy(by FindPointTransactionBy) entity.PointTransaction {
	pt := entity.PointTransaction{}
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
func (repo PointTransactionRepository) GetPageListBy(by GetPointTransactionPageListBy) ([]entity.PointTransaction, int64) {
	list := make([]entity.PointTransaction, 0)

	db := repo.ctx.DB.Model(entity.PointTransaction{})
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

func (repo PointTransactionRepository) CountByToday(by GetPointTransactionCountBy) ([]map[string]interface{}, int64, error) {
	var result []map[string]interface{}
	db := repo.ctx.DB.Model(&entity.PointTransaction{})
	if len(by.OpenIds) > 0 {
		db.Where("openid in (?)", by.OpenIds)
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
	db.Where("date(create_time) = CURRENT_DATE")
	var count int64
	if err := db.Count(&count).Find(&result).Error; err != nil {
		return result, count, err
	}
	return result, count, nil
}
