package event

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/event"
	"mio/internal/pkg/repository/repotypes"
)

var DefaultEventRepository = EventRepository{DB: app.DB}

type EventRepository struct {
	DB *gorm.DB
}

func (repo EventRepository) Save(ev *event.Event) error {
	return repo.DB.Save(ev).Error
}
func (repo EventRepository) Create(ev *event.Event) error {
	return repo.DB.Create(ev).Error
}
func (repo EventRepository) FindEvent(by FindEventBy) (event.Event, error) {
	ev := event.Event{}
	db := repo.DB.Model(ev)

	if by.ProductItemId != "" {
		db.Where("product_item_id = ?", by.ProductItemId)
	}
	if by.EventId != "" {
		db.Where("event_id = ?", by.EventId)
	}
	if by.Active.Valid {
		db.Where("active = ?", by.Active.Bool)
	}

	err := db.Take(&ev).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return event.Event{}, err
	}

	return ev, nil
}

func (repo EventRepository) FindEventCate(by FindEventBy) (repotypes.EventRet, error) {
	ev := event.Event{}
	var Ret repotypes.EventRet

	db := repo.DB.Model(ev)
	db.Select("event.title as title,event_category.title as cateTitle")
	if by.ProductItemId != "" {
		db.Where("product_item_id = ?", by.ProductItemId)
	}
	db.Joins("left join event_category on event_category.event_category_id = \"event\".event_category_id ")
	err := db.First(&Ret).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return Ret, err
	}
	return Ret, nil
}

func (repo EventRepository) GetEventPageList(by GetEventPageListBy) (list []event.Event, total int64, err error) {

	list = make([]event.Event, 0)

	db := repo.DB.Model(event.Event{})
	for _, orderBy := range by.OrderBy {
		switch orderBy {
		case event.OrderByEventSortDesc:
			db.Where("sort desc")
		}
	}

	if by.EventCategoryId != "" {
		db.Where("event_category_id = ?", by.EventCategoryId)
	}

	if by.IsShow != 0 {
		db.Where("is_show = ?", by.IsShow)
	}

	err = db.Count(&total).Offset(by.Offset).Limit(by.Limit).Find(&list).Error

	return
}
func (repo EventRepository) GetEventList(by GetEventListBy) (list []event.Event, err error) {

	list = make([]event.Event, 0)

	db := repo.DB.Model(event.Event{})
	for _, orderBy := range by.OrderBy {
		switch orderBy {
		case event.OrderByEventSortDesc:
			db.Order("sort desc")
		}
	}

	if by.EventCategoryId != "" {
		db.Where("event_category_id = ?", by.EventCategoryId)
	}
	if by.Active.Valid {
		db.Where("active = ?", by.Active.Bool)
	}

	if by.IsShow != 0 {
		db.Where("is_show = ?", by.IsShow)
	}

	err = db.Find(&list).Error

	return
}