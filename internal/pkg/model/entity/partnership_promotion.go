package entity

import "mio/internal/pkg/model"

type PartnershipPromotionTrigger string

const (
	PartnershipPromotionTriggerREGISTER  PartnershipPromotionTrigger = "REGISTER"
	PartnershipPromotionTriggerBINDPHONE PartnershipPromotionTrigger = "BIND_PHONE"
	PartnershipPromotionTriggerOCR       PartnershipPromotionTrigger = "OCR"
)

type PartnershipPromotion struct {
	ID           int64                       `json:"id"`
	PromotionId  string                      `json:"promotionId"`
	Partnership  string                      `json:"partnership"`
	StartTime    model.Time                  `json:"startTime"`
	EndTime      model.Time                  `json:"endTime"`
	CouponTypeId string                      `json:"couponTypeId"`
	Trigger      PartnershipPromotionTrigger `json:"trigger"`
}
