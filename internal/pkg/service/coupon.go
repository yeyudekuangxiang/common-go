package service

import (
	"mio/internal/pkg/model/entity/coupon"
	"mio/internal/pkg/repository"
	"strings"
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

func (r CouponService) CouponListOfOpenid(openid string, couponTypeIds []string) ([]Coupon.CouponRes, error) {
	res, err := r.r.CouponListOfOpenid(openid, couponTypeIds)
	var res2 []Coupon.CouponRes
	if err == nil {
		for _, row := range res {
			row.CouponType = "兑换券"
			timeArr1 := strings.Split(row.UpdateTime, "T")
			time1 := timeArr1[0]
			timeArr2 := strings.Split(timeArr1[1], ".")
			time2 := timeArr2[0]
			row.UpdateTime = time1 + " " + time2
			res2 = append(res2, row)
		}
	}
	return res2, err
}
