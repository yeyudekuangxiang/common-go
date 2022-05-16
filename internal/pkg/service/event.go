package service

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultEventService = EventService{}

type EventService struct {
}

func (srv EventService) FindEvent(by FindEventBy) (*entity.Event, error) {
	event := entity.Event{}
	db := app.DB.Model(event)

	if by.ProductItemId != "" {
		db.Where("product_item_id = ?", by.ProductItemId)
	}

	err := db.Take(&event).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	return &event, nil
}
