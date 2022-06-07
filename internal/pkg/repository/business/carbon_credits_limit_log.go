package business

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/business"
)

var DefaultCarbonCreditsLimitLogRepository = CarbonCreditsLimitLogRepository{DB: app.DB}

type CarbonCreditsLimitLogRepository struct {
	DB *gorm.DB
}

func (repo CarbonCreditsLimitLogRepository) Save(log *business.CarbonCreditsLimitLog) error {
	return repo.DB.Save(log).Error
}

func (repo CarbonCreditsLimitLogRepository) Create(log *business.CarbonCreditsLimitLog) error {
	return repo.DB.Create(log).Error
}
func (repo CarbonCreditsLimitLogRepository) FindLimitLog(by FindCarbonCreditsLimitLogBy) business.CarbonCreditsLimitLog {
	log := business.CarbonCreditsLimitLog{}
	db := repo.DB.Model(log)

	if !by.TimePoint.IsZero() {
		db.Where("time_point = ?", by.TimePoint)
	}
	if by.Type != "" {
		db.Where("type = ?", by.Type)
	}
	if by.UserId != 0 {
		db.Where("b_user_id = ?", by.UserId)
	}

	err := db.Take(&log).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return log
}
