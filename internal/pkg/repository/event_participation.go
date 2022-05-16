package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultEventParticipationRepository = EventParticipationRepository{DB: app.DB}

type EventParticipationRepository struct {
	DB *gorm.DB
}

func (repo EventParticipationRepository) CreateBatch(list *[]entity.EventParticipation) error {
	return repo.DB.Create(list).Error
}
func (repo EventParticipationRepository) Create(participation *entity.EventParticipation) error {
	return repo.DB.Create(participation).Error
}
