package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity/coupon"
)

type RedeemCodeRepository struct {
	ctx *context.MioContext
}

func NewRedeemCodeRepository(ctx *context.MioContext) *RedeemCodeRepository {
	return &RedeemCodeRepository{ctx: ctx}
}

func (repo *RedeemCodeRepository) Save(code *coupon.RedeemCode) error {
	return repo.ctx.DB.Save(code).Error
}
func (repo *RedeemCodeRepository) Create(code *coupon.RedeemCode) error {
	return repo.ctx.DB.Create(code).Error
}

func (repo *RedeemCodeRepository) GetRedeemCode(by GetRedeemCodeBy) (*coupon.RedeemCode, bool, error) {
	code := coupon.RedeemCode{}
	db := repo.ctx.DB.Model(code)

	if by.CouponId != "" {
		db.Where("coupon_id = ?", by.CouponId)
	}

	if by.CodeId != "" {
		db.Where("code_id = ?", by.CodeId)
	}

	err := db.Take(&code).Error

	if err == nil {
		return &code, true, nil
	}
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	return nil, false, err
}
