package service

import (
	"mio/internal/pkg/model/entity/coupon"
	"mio/internal/pkg/repository"
)

var DefaultCouponTypeService = CouponTypeService{repo: repository.DefaultCouponTypeRepository}

type CouponTypeService struct {
	repo repository.CouponTypeRepository
}

func (srv CouponTypeService) FindCouponType(by FindCouponTypeBy) (*coupon.CouponType, error) {
	ct := srv.repo.FindCouponType(repository.FindCouponTypeBy{
		CouponTypeId: by.CouponTypeId,
	})
	return &ct, nil
}
