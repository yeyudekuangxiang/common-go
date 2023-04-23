package service

import (
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/pkg/errno"
	"time"
)

var DefaultPartnershipRedemptionService = PartnershipRedemptionService{}

type PartnershipRedemptionService struct {
}

// ProcessPromotionInformation 第三方平台优惠券兑换
func (srv PartnershipRedemptionService) ProcessPromotionInformation(openId string, partnership entity.PartnershipType, trigger entity.PartnershipPromotionTrigger) ([]entity.PartnershipRedemption, error) {
	app.Logger.Info("根据第三方获取合作活动列表", openId, partnership, trigger)
	promotionList, err := DefaultPartnershipPromotionService.GetPartnershipPromotionList(GetPartnershipPromotionListBy{
		Partnership: partnership,
		Trigger:     trigger,
	})

	if err != nil {
		return nil, err
	}

	if len(promotionList) == 0 {
		return nil, nil
	}

	infoList, err := srv.formatInfoList(openId, promotionList)
	if err != nil {
		return nil, err
	}

	return infoList, app.DB.Create(&infoList).Error
}
func (srv PartnershipRedemptionService) formatInfoList(openId string, promotionList []entity.PartnershipPromotion) ([]entity.PartnershipRedemption, error) {
	t := time.Now()
	infoList := make([]entity.PartnershipRedemption, 0)
	app.Logger.Infof("循环第三方活动列表 %+v", promotionList)
	for _, promotion := range promotionList {
		app.Logger.Info("查询活动相关的优惠券类型信息", promotion, promotion.CouponTypeId)
		//查询活动相关的优惠券模版信息
		coupon, err := DefaultCouponTypeService.FindCouponType(FindCouponTypeBy{
			CouponTypeId: promotion.CouponTypeId,
		})
		if err != nil {
			app.Logger.Error("查询优惠券类型失败", promotion, err)
			return nil, err
		}

		if coupon.ID == 0 || !coupon.Active {
			return nil, err
		}

		if coupon.External {
			item, err := srv.findCouponToSend(openId, promotion, t)
			if err != nil {
				return nil, err
			}
			infoList = append(infoList, *item)
		} else {
			itemList, err := srv.createCouponToSend(openId, promotion, t)
			if err != nil {
				return nil, err
			}
			infoList = append(infoList, itemList...)
		}
	}
	return infoList, nil
}
func (srv PartnershipRedemptionService) findCouponToSend(openId string, promotion entity.PartnershipPromotion, now time.Time) (*entity.PartnershipRedemption, error) {
	app.Logger.Info("获取一张未发放的优惠券", promotion.ID, promotion.CouponTypeId)
	//获取一张未发放的优惠券
	unSendCoupon, err := DefaultCouponService.RetrieveUnassignedCoupon(promotion.CouponTypeId)
	if err != nil {
		app.Logger.Error(promotion, err)
		return nil, err
	}
	if unSendCoupon.ID == 0 {
		return nil, errno.ErrCommon.WithMessage("兑换券存量不足")
	}

	app.Logger.Info("获取一张未发放的优惠券", openId, unSendCoupon.CouponTypeId, entity.OrderTypeRedeem)
	//兑换优惠券
	_, err = DefaultCouponService.RedeemCoupon(RedeemCouponParam{
		OpenId:    openId,
		CouponId:  unSendCoupon.CouponId,
		OrderType: entity.OrderTypeRedeem,
	})

	if err != nil {
		app.Logger.Error(promotion, err)
		return nil, err
	}

	return &entity.PartnershipRedemption{
		OpenId:      openId,
		Time:        model.Time{Time: now},
		CouponId:    model.NullString(unSendCoupon.CouponId),
		PromotionId: model.NullString(promotion.PromotionId),
	}, nil
}
func (srv PartnershipRedemptionService) createCouponToSend(openId string, promotion entity.PartnershipPromotion, now time.Time) ([]entity.PartnershipRedemption, error) {
	//生成优惠券发放
	couponIds, err := DefaultCouponService.GenerateCouponBatch(GenerateCouponBatchParam{
		CouponTypeId: promotion.CouponTypeId,
		BatchSize:    1,
	})
	if err != nil {
		return nil, err
	}
	app.Logger.Info("生成优惠券用于发放", promotion.ID, promotion.CouponTypeId, couponIds)

	infoList := make([]entity.PartnershipRedemption, 0)
	for _, couponId := range couponIds {
		app.Logger.Info("兑换优惠券", openId, couponId, entity.OrderTypeRedeem)
		_, err := DefaultCouponService.RedeemCoupon(RedeemCouponParam{
			OpenId:    openId,
			CouponId:  couponId,
			OrderType: entity.OrderTypeRedeem,
		})
		if err != nil {
			return nil, err
		}
		infoList = append(infoList, entity.PartnershipRedemption{
			OpenId:      openId,
			Time:        model.Time{Time: now},
			CouponId:    model.NullString(couponId),
			PromotionId: model.NullString(promotion.PromotionId),
		})
	}
	return infoList, nil
}
