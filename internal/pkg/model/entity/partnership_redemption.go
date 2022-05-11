package entity

import "mio/internal/pkg/model"

type PartnershipType string

const (
	PartnershipTypeXINGXINGCHONGDIAN PartnershipType = "XINGXINGCHONGDIAN"
	PartnershipTypeEVCARD            PartnershipType = "EVCARD"
	PartnershipTypeLIANTONG          PartnershipType = "LIANTONG"
	PartnershipTypeDIDI              PartnershipType = "DIDI"
)

type PartnershipRedemption struct {
	ID          int64            `json:"id"`
	Time        model.Time       `json:"time"`
	OpenId      string           `json:"openId" gorm:"column:'openid'"`
	PromotionId model.NullString `json:"promotionId" gorm:"column:'promotion_id'"`
	CouponId    model.NullString `json:"couponId" gorm:"column:'coupon_id'"`
}
