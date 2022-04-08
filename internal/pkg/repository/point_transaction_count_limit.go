package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultPointTransactionCountLimitRepository = PointTransactionCountLimitRepository{DB: app.DB}

type PointTransactionCountLimitRepository struct {
	DB *gorm.DB
}

func (p PointTransactionCountLimitRepository) Save(transactionCountLimit *entity.PointTransactionCountLimit) error {
	return p.DB.Save(transactionCountLimit).Error
}

func (p PointTransactionCountLimitRepository) FindBy(by FindPointTransactionCountLimitBy) entity.PointTransactionCountLimit {
	limit := entity.PointTransactionCountLimit{}

	db := p.DB.Model(entity.PointTransactionCountLimit{})
	if by.OpenId != "" {
		db.Where("openid = ?", by.OpenId)
	}
	if by.TransactionType != "" {
		db.Where("transaction_type = ?", by.TransactionType)
	}
	if !by.TransactionDate.IsZero() {
		db.Where("transaction_date = ?", by.TransactionDate)
	}
	if err := db.First(&limit).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return limit
}
