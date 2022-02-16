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

func (r CouponService) CouponListOfOpenid(openid string, couponTypeIds string) ([]Coupon.CouponRes, error) {
	res, err := r.r.CouponListOfOpenid(openid, couponTypeIds)
	var res2 []Coupon.CouponRes
	if err == nil {
		for _, row := range res {
			row.CouponType = "滴滴兑换券"
			res2 = append(res2, row)
		}
	}
	return res2, err
}
