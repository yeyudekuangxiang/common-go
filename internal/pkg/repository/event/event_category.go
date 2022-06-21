package event

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/event"
)

var DefaultEventCategoryRepository = EventCategoryRepository{DB: app.DB}

type EventCategoryRepository struct {
	DB *gorm.DB
}

func (repo EventCategoryRepository) GetEventCategoryList(by GetEventCategoryListBy) ([]event.EventCategory, error) {
	list := make([]event.EventCategory, 0)
	db := repo.DB.Model(event.EventCategory{})
	for _, orderBy := range by.OrderBy {
		switch orderBy {
		case event.OrderByEventCategorySortDesc:
			db.Order("sort desc")
		}
	}
	return list, db.Find(&list).Error
}
