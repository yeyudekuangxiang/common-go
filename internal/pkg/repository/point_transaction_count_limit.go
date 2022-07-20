package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
)

type PointTransactionCountLimitRepository struct {
	ctx *context.MioContext
}

func NewPointTransactionCountLimitRepository(ctx *context.MioContext) *PointTransactionCountLimitRepository {
	return &PointTransactionCountLimitRepository{ctx: ctx}
}

func (p PointTransactionCountLimitRepository) Save(transactionCountLimit *entity.PointTransactionCountLimit) error {
	return p.ctx.DB.Save(transactionCountLimit).Error
}

func (p PointTransactionCountLimitRepository) FindBy(by FindPointTransactionCountLimitBy) entity.PointTransactionCountLimit {
	limit := entity.PointTransactionCountLimit{}

	db := p.ctx.DB.Model(entity.PointTransactionCountLimit{})
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
