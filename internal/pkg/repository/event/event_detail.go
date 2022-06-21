package event

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/event"
)

var DefaultEventDetailRepository = EventDetailRepository{DB: app.DB}

type EventDetailRepository struct {
	DB *gorm.DB
}

func (repo EventDetailRepository) GetEventDetailList(by GetEventDetailListBy) ([]event.EventDetail, error) {
	list := make([]event.EventDetail, 0)

	db := repo.DB.Model(event.EventDetail{})

	if by.EventId != "" {
		db.Where("event_id = ?", by.EventId)
	}
	return list, db.Order("id asc").Find(&list).Error
}
