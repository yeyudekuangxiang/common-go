package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/coupon"
)

var DefaultCouponTypeRepository = CouponTypeRepository{DB: app.DB}

type CouponTypeRepository struct {
	DB *gorm.DB
}

func (repo CouponTypeRepository) FindCouponType(by FindCouponTypeBy) coupon.CouponType {
	ct := coupon.CouponType{}
	db := repo.DB.Model(ct)

	if by.CouponTypeId != "" {
		db.Where("coupon_type_id = ?", by.CouponTypeId)
	}

	err := db.First(&ct).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return ct
}
