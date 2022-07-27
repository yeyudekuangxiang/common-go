package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultCheckinHistoryRepository = CheckinHistoryRepository{DB: app.DB}

type CheckinHistoryRepository struct {
	DB *gorm.DB
}

func (repo CheckinHistoryRepository) FindCheckinHistory(by FindCheckinHistoryBy) entity.CheckinHistory {
	history := entity.CheckinHistory{}
	db := repo.DB.Model(history)
	if by.OpenId != "" {
		db.Where("openid = ?", by.OpenId)
	}
	repo.orderBy(db, by.OrderBy)

	err := db.Take(&history).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return history
}
func (repo CheckinHistoryRepository) orderBy(db *gorm.DB, list entity.OrderByList) {
	for _, orderBy := range list {
		switch orderBy {
		case entity.OrderByCheckinHistoryTimeDesc:
			db.Order("time desc")
		}
	}
}
func (repo CheckinHistoryRepository) Create(history *entity.CheckinHistory) error {
	return repo.DB.Create(history).Error
}
