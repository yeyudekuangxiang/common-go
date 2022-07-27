package event

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/event"
)

var DefaultEventRuleRepository = EventRuleRepository{DB: app.DB}

type EventRuleRepository struct {
	DB *gorm.DB
}

func (repo EventRuleRepository) GetEventRuleList(by GetEventRuleListBy) ([]event.EventRule, error) {
	list := make([]event.EventRule, 0)

	db := repo.DB.Model(event.EventRule{})

	if by.EventId != "" {
		db.Where("event_id = ?", by.EventId)
	}
	return list, db.Order("id asc").Find(&list).Error
}
