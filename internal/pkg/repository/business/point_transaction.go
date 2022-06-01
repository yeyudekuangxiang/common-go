package business

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/business"
)

var DefaultPointTransactionRepository = NewPointTransactionRepository(app.DB)

func NewPointTransactionRepository(db *gorm.DB) PointTransactionRepository {
	return PointTransactionRepository{
		DB: db,
	}
}

type PointTransactionRepository struct {
	DB *gorm.DB
}

func (p PointTransactionRepository) Save(transaction *business.PointTransaction) error {
	return p.DB.Save(transaction).Error
}

func (p PointTransactionRepository) GetListBy(by GetPointTransactionListBy) []business.PointTransaction {
	list := make([]business.PointTransaction, 0)

	db := p.DB.Model(business.PointTransaction{})
	if by.UserId != 0 {
		db.Where("b_user_id = ?", by.UserId)
	}

	if !by.StartTime.IsZero() {
		db.Where("created_at >= ?", by.StartTime)
	}
	if !by.EndTime.IsZero() {
		db.Where("updated_at <= ?", by.EndTime)
	}

	if by.Type != "" {
		db.Where("type = ?", by.Type)
	}

	for _, orderBy := range by.OrderBy {
		switch orderBy {
		case business.OrderByPointTranCTDESC:
			db.Order("created_at desc")
		}
	}

	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}

	return list
}
