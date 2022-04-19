package service

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultDuiBaActivityService = DuiBaActivityService{}

type DuiBaActivityService struct {
}

func (srv DuiBaActivityService) FindActivity(activityId string) (*entity.DuiBaActivity, error) {
	activity := entity.DuiBaActivity{}
	err := app.DB.Where("activity_id = ?", activityId).First(&activity).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	return &activity, nil
}
