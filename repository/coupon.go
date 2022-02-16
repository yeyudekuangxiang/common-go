package repository

import (
	"mio/core/app"
	Coupon "mio/model/Coupon"
)

var DefaultCouponRepository ICouponRepository = NewCouponRepository()

type ICouponRepository interface {
	CouponListOfOpenid(openid string, couponTypeIds string) ([]Coupon.Coupon, error)
}

func NewCouponRepository() CouponRepository {
	return CouponRepository{}
}

type CouponRepository struct {
}

func (p CouponRepository) CouponListOfOpenid(openid string, couponTypeIds string) ([]Coupon.Coupon, error) {
	var Coupons []Coupon.Coupon
	if err := app.DB.Table("coupon").Where("Coupon_type_id in (?)", couponTypeIds).Where("openid = ?", openid).Where("redeemed = ?", "true").Find(&Coupons).Error; err != nil {
		return nil, err
	}
	return Coupons, nil
}
