package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultStepHistoryRepository = StepHistoryRepository{DB: app.DB}

type StepHistoryRepository struct {
	DB *gorm.DB
}

func (repo StepHistoryRepository) FindBy(by FindStepHistoryBy) entity.StepHistory {
	sh := entity.StepHistory{}
	db := repo.DB.Model(sh)

	if by.OpenId != "" {
		db.Where("openid = ?", by.OpenId)
	}
	if !by.Day.IsZero() {
		db.Where("recorded_time = ?", by.Day)
	}

	for _, orderBy := range by.OrderBy {
		switch orderBy {
		case entity.OrderByStepHistoryCountDesc:
			db.Order("count desc")
		}
	}

	err := db.First(&sh).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	return sh
}
