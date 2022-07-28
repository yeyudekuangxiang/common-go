package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultCarbonTransactionRepository = NewCarbonTransactionRepository(app.DB)

func NewCarbonTransactionRepository(db *gorm.DB) CarbonTransactionRepository {
	return CarbonTransactionRepository{
		DB: db,
	}
}

type CarbonTransactionRepository struct {
	DB *gorm.DB
}

func (p CarbonTransactionRepository) Save(transaction *entity.CarbonTransaction) error {
	return p.DB.Save(transaction).Error
}

func (repo CarbonTransactionRepository) Create(transaction *entity.CarbonTransaction) error {
	return repo.DB.Create(transaction).Error
}

func (p CarbonTransactionRepository) GetListBy(by GetPointTransactionListBy) []entity.CarbonTransaction {
	list := make([]entity.CarbonTransaction, 0)

	db := p.DB.Model(entity.CarbonTransaction{})
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
func (p CarbonTransactionRepository) FindBy(by FindPointTransactionBy) entity.CarbonTransaction {
	pt := entity.CarbonTransaction{}
	db := p.DB.Model(pt)
	if by.TransactionId != "" {
		db.Where("transaction_id", by.TransactionId)
	}
	err := db.First(&pt).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return pt
}
func (p CarbonTransactionRepository) GetPageListBy(by GetPointTransactionPageListBy) ([]entity.CarbonTransaction, int64) {
	list := make([]entity.CarbonTransaction, 0)

	db := p.DB.Model(entity.CarbonTransaction{})
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
