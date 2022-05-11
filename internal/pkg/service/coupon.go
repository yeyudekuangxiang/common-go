package service

import (
	"errors"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
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

func (r CouponService) CouponListOfOpenid(openid string, couponTypeIds []string) ([]coupon.CouponRes, error) {
	res, err := r.r.CouponListOfOpenid(openid, couponTypeIds)
	var res2 []coupon.CouponRes
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

func (r CouponService) FindCoupon(by FindCouponBy) (*coupon.Coupon, error) {
	cp := r.r.FindCoupon(repository.FindCouponBy{
		CouponTypeId: by.CouponTypeId,
	})
	return &cp, nil
}

//RetrieveUnassignedCoupon 获取未发放的优惠券
func (r CouponService) RetrieveUnassignedCoupon(couponTypeId string) (*coupon.Coupon, error) {
	cp := coupon.Coupon{}

	err := app.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(cp).Set("gorm:query_option", "for update skip locked").
			Where("coupon_type_id = ?", couponTypeId).
			Where("update_time is null").
			Where("openid is null").Take(&cp).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			panic(err)
		}

		if cp.ID == 0 {
			app.Logger.Error("兑换券存量不足", couponTypeId)
			return errors.New("兑换券存量不足")
		}
		cp.UpdateTime = model.NewTime()
		err = tx.Save(&cp).Error
		return err
	})
	return &cp, err
}

// RedeemCoupon 兑换优惠券
func (r CouponService) RedeemCoupon(param RedeemCouponParam) {
	/*coupon := r.r.FindCoupon(repository.FindCouponBy{
		CouponId: param.CouponId,
	})
	return*/
}
