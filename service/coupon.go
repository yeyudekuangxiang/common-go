package service

import (
	Coupon "mio/model/Coupon"
	"mio/repository"
)

var DefaultCouponService = NewCouponService(repository.DefaultCouponRepository)

func NewCouponService(r repository.ICouponRepository) CouponService {
	return CouponService{
		r: r,
	}
}

type CouponService struct {
	r repository.ICouponRepository
}

func (r CouponService) CouponListOfOpenid(openid string, couponTypeIds string) ([]Coupon.Coupon, error) {
	return r.r.CouponListOfOpenid(openid, couponTypeIds)
}
