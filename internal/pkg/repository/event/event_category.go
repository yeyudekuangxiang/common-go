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

func (repo EventCategoryRepository) GetEventCategoryList(by GetEventCategoryListBy) ([]event.APIEventCategory, error) {
	db := repo.DB.Model(event.EventCategory{}).Preload("Link", func(db *gorm.DB) *gorm.DB {
		if by.Display >= 0 {
			return db.Where("display in ?", []int{0, by.Display})
		}
		return db
	})

	for _, orderBy := range by.OrderBy {
		switch orderBy {
		case event.OrderByEventCategorySortDesc:
			db.Order("sort desc")
		}
	}

	if by.Active.Valid {
		db.Where("active = ?", by.Active.Bool)
	}

	list := make([]event.APIEventCategory, 0)
	return list, db.Find(&list).Error
}
