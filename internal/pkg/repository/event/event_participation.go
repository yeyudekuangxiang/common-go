package event

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/event"
)

var DefaultEventParticipationRepository = EventParticipationRepository{DB: app.DB}

type EventParticipationRepository struct {
	DB *gorm.DB
}

func (repo EventParticipationRepository) CreateBatch(list *[]event.EventParticipation) error {
	return repo.DB.Create(list).Error
}
func (repo EventParticipationRepository) Create(participation *event.EventParticipation) error {
	return repo.DB.Create(participation).Error
}
func (repo EventParticipationRepository) GetParticipationPageList(by GetParticipationPageListBy) ([]event.EventParticipation, int64, error) {
	list := make([]event.EventParticipation, 0)
	var total int64
	db := repo.DB.Model(event.EventParticipation{})

	for _, orderBy := range by.OrderBy {
		switch orderBy {
		case event.OrderByEventParticipationCountDesc:
			db.Order("count desc")
		case event.OrderByEventParticipationTimeDesc:
			db.Order("time desc")
		}
	}
	if by.EventId != "" {
		db.Where("event_id = ?", by.EventId)
	}
	return list, total, db.Count(&total).Limit(by.Limit).Offset(by.Offset).Find(&list).Error
}
