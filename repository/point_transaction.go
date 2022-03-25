package repository

import (
	"gorm.io/gorm"
	"mio/core/app"
	"mio/model/entity"
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

func (p PointTransactionRepository) Save(transaction *entity.PointTransaction) error {
	return p.DB.Save(transaction).Error
}

func (p PointTransactionRepository) GetListBy(by GetPointTransactionListBy) []entity.PointTransaction {
	list := make([]entity.PointTransaction, 0)

	db := p.DB.Model(entity.PointTransaction{})
	if by.OpenId != "" {
		db.Where("openid = ?", by.OpenId)
	}
	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return list
}
