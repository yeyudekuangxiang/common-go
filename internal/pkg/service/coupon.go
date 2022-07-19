package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/model/entity/coupon"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/repository/repotypes"
	"mio/internal/pkg/service/product"
	"mio/internal/pkg/service/srv_types"

	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"strings"
	"time"
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
			//return errors.New("兑换券存量不足")
			return nil
		}
		cp.UpdateTime = model.NewTime()
		err = tx.Save(&cp).Error
		return err
	})
	return &cp, err
}

// RedeemCoupon 兑换优惠券
func (r CouponService) RedeemCoupon(param RedeemCouponParam) (*RedeemCouponWithTransactionResult, error) {
	if param.OrderType == "" {
		param.OrderType = entity.OrderTypeRedeem
	} else if param.OrderType == entity.OrderTypeGreenTorch {
		return r.RedeemCouponWithTransaction(RedeemCouponWithTransactionParam{
			OpenId:               param.OpenId,
			CouponId:             param.CouponId,
			OrderType:            param.OrderType,
			PointTransactionType: entity.POINT_GREEN_TORCH,
		})
	}
	return r.RedeemCouponWithTransaction(RedeemCouponWithTransactionParam{
		OpenId:               param.OpenId,
		CouponId:             param.CouponId,
		OrderType:            param.OrderType,
		PointTransactionType: entity.POINT_COUPON,
	})
}
func (r CouponService) RedeemCouponWithTransaction(param RedeemCouponWithTransactionParam) (*RedeemCouponWithTransactionResult, error) {
	coupon := r.r.FindCoupon(repository.FindCouponBy{
		CouponId: param.CouponId,
	})
	if coupon.ID == 0 {
		return nil, errors.New("券码不存在")
	}

	if !r.IsActiveCoupon(coupon.StartTime.Time, coupon.EndTime.Time) {
		return nil, errors.New("不在券码有效期内")
	}

	if coupon.Redeemed {
		return nil, errors.New("券码已被兑换")
	}

	coupon.Openid = param.OpenId
	coupon.UpdateTime = model.NewTime()
	coupon.Redeemed = true

	contentType, err := DefaultCouponTypeService.FindCouponType(FindCouponTypeBy{
		CouponTypeId: coupon.CouponTypeId,
	})
	if err != nil {
		return nil, err
	}
	if contentType.ID == 0 {
		return nil, errors.New("未找到券码类型")
	}

	orderId := ""
	if contentType.Point != 0 {
		err := r.RedeemCouponToPoints(param.OpenId, int(contentType.Point), coupon.CouponId, param.PointTransactionType)
		if err != nil {
			return nil, err
		}
	} else if contentType.ProductItemId != "" {
		order, err := r.RedeemCouponToItems(param.OpenId, param.OrderType, *contentType, entity.PartnershipType(contentType.Partnership))
		if err != nil {
			return nil, err
		}
		orderId = order.OrderId
	} else {
		return nil, errors.New("改兑换码为第三方伙伴的券码哦～请去对应的地方兑换")
	}

	err = r.r.Save(&coupon)
	if err != nil {
		return nil, err
	}
	return &RedeemCouponWithTransactionResult{Point: int(contentType.Point), OrderId: orderId}, nil
}

// IsActiveCoupon 判断优惠券是否在有效期
func (r CouponService) IsActiveCoupon(startTime, endTime time.Time) bool {
	now := time.Now()
	if !startTime.IsZero() && startTime.After(now) {
		return false
	}

	if !endTime.IsZero() && endTime.Before(now) {
		return false
	}

	return true
}

// RedeemCouponToPoints 将优惠券兑换成积分
func (r CouponService) RedeemCouponToPoints(openId string, value int, couponId string, pt entity.PointTransactionType) error {
	_, err := DefaultPointTransactionService.Create(CreatePointTransactionParam{
		OpenId:       openId,
		Value:        value,
		Type:         pt,
		AdditionInfo: fmt.Sprintf("Consume the coupon %s", couponId),
	})
	return err
}

// RedeemCouponToItems 将优惠券兑换成订单
func (r CouponService) RedeemCouponToItems(openId string, orderType entity.OrderType, couponType coupon.CouponType, partnership entity.PartnershipType) (*entity.Order, error) {
	user, err := DefaultUserService.GetUserByOpenId(openId)

	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, errno.ErrUserNotFound
	}

	if couponType.ProductItemId == "" {
		return nil, errors.New("invalid product item")
	}
	productItem, err := product.DefaultProductItemService.FindProductByItemId(string(couponType.ProductItemId))
	if err != nil {
		return nil, err
	}
	if productItem.ID == 0 {
		return nil, errors.New("invalid product item")
	}

	address, err := DefaultAddressService.FindDefaultAddress(user.OpenId)
	if err != nil {
		return nil, err
	}

	return DefaultOrderService.SubmitOrder(SubmitOrderParam{
		Order: SubmitOrder{
			UserId:    user.ID,
			AddressId: address.AddressId,
			OrderType: orderType,
		},
		Items: []SubmitOrderItem{
			{
				ItemId: string(couponType.ProductItemId),
				Count:  int(couponType.ProductItemCount),
			},
		},
		PartnershipType: partnership,
	})
}

// GenerateCouponBatch 生成优惠券
func (r CouponService) GenerateCouponBatch(param GenerateCouponBatchParam) ([]string, error) {
	contentType, err := DefaultCouponTypeService.FindCouponType(FindCouponTypeBy{
		CouponTypeId: param.CouponTypeId,
	})
	if err != nil {
		return nil, err
	}

	couponIds := make([]string, 0)
	cps := make([]coupon.Coupon, 0)
	for i := 0; i < param.BatchSize; i++ {
		id := util.UUID()
		cp := coupon.Coupon{
			CouponId:     id,
			CouponTypeId: contentType.CouponTypeId,
			CreateTime:   model.NewTime(),
			StartTime:    contentType.StartTime,
			EndTime:      contentType.EndTime,
		}
		couponIds = append(couponIds, id)
		cps = append(cps, cp)
	}

	return couponIds, r.r.CreateBatch(&cps)
}
func (r CouponService) GetPageUserCouponRecord(getCouponDTO srv_types.GetPageCouponRecordDTO) ([]srv_types.BaseCouponRecordDTO, int64, error) {
	getCouponDO := repotypes.GetPageUserCouponTypeDO{}
	if err := util.MapTo(getCouponDTO, &getCouponDO); err != nil {
		return nil, 0, err
	}

	couponTypeList, total, err := r.r.GetPageCouponRecord(getCouponDO)
	if err != nil {
		return nil, 0, err
	}

	couponRecordList := make([]srv_types.BaseCouponRecordDTO, 0)

	for _, couponType := range couponTypeList {
		coverImage := util.LinkJoin(config.Config.OSS.CdnDomain, "/static/mp2c/images/coupon/redeem-point-icon.png")
		if couponType.ProductItemId != "" {
			productItem, err := product.DefaultProductItemService.FindProductByItemId(couponType.ProductItemId)
			if err != nil {
				return nil, 0, err
			}
			coverImage = productItem.ImageUrl
		}

		couponRecordList = append(couponRecordList, srv_types.BaseCouponRecordDTO{
			ID:         couponType.ID,
			CoverImage: coverImage,
			Title:      couponType.Name,
			UpdateDate: couponType.UpdateTime.Date(),
		})
	}
	return couponRecordList, total, nil
}
