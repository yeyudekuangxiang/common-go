package service

import (
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultPartnershipPromotionService = PartnershipPromotionService{}

type PartnershipPromotionService struct {
}

func (srv PartnershipPromotionService) GetPromotionPromotionList(by GetPromotionPromotionListBy) ([]entity.PartnershipPromotion, error) {
	list := make([]entity.PartnershipPromotion, 0)
	db := app.DB.Model(entity.PartnershipPromotion{})

	if by.Partnership != "" {
		db.Where("partnership = ?", by.Partnership)
	}

	if by.Trigger != "" {
		db.Where("trigger = ?", by.Trigger)
	}
	err := db.Find(&list).Error
	if err != nil {
		panic(err)
	}

	return list, nil
}
