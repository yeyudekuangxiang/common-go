package service

import (
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

type PartnershipRedemption struct {
}

//第三方平台优惠券兑换
func (srv PartnershipRedemption) ProcessPromotionInformation(openId string, partnership entity.PartnershipType, trigger entity.PartnershipPromotionTrigger) error {
	promotionList, err := DefaultPartnershipPromotionService.GetPromotionPromotionList(GetPromotionPromotionListBy{
		partnership,
		trigger,
	})

	if err != nil {
		return err
	}

	for _, promotion := range promotionList {
		coupon, err := DefaultCouponTypeService.FindCouponType(FindCouponTypeBy{
			CouponTypeId: promotion.CouponTypeId,
		})
		if err != nil {
			app.Logger.Error(promotion, err)
			continue
		}

		if coupon.ID == 0 || !coupon.Active {
			continue
		}

		if coupon.External {
			/*unSendCoupon, err := DefaultCouponService.RetrieveUnassignedCoupon(promotion.CouponTypeId)
			if err != nil {
				app.Logger.Error(promotion, err)
				continue
			}*/

		}

	}
	return nil
}
