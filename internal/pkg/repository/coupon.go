package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/coupon"
)

var DefaultCouponRepository ICouponRepository = NewCouponRepository(app.DB)

type ICouponRepository interface {
	CouponListOfOpenid(openid string, couponTypeIds []string) ([]coupon.CouponRes, error)
	FindCoupon(by FindCouponBy) coupon.Coupon
}

func NewCouponRepository(DB *gorm.DB) CouponRepository {
	return CouponRepository{DB: DB}
}

type CouponRepository struct {
	DB *gorm.DB
}

func (p CouponRepository) CouponListOfOpenid(openid string, couponTypeIds []string) ([]coupon.CouponRes, error) {
	var Coupons []coupon.CouponRes
	if err := p.DB.Table("coupon").Where("Coupon_type_id in (?)", couponTypeIds).Where("openid = ?", openid).Where("redeemed = ?", "true").Find(&Coupons).Error; err != nil {
		return nil, err
	}
	return Coupons, nil
}
func (p CouponRepository) FindCoupon(by FindCouponBy) coupon.Coupon {
	cp := coupon.Coupon{}
	db := p.DB.Model(cp)

	if by.CouponTypeId != "" {
		db.Where("coupon_type_id = ?", by.CouponTypeId)
	}
	if by.CouponId != "" {
		db.Where("coupon_id = ?", by.CouponId)
	}

	err := db.First(&cp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return cp
}
