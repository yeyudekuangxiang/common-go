package business

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/business"
)

var DefaultPointLimitLogRepository = PointLimitLogRepository{DB: app.DB}

type PointLimitLogRepository struct {
	DB *gorm.DB
}

func (repo PointLimitLogRepository) Save(log *business.PointLimitLog) error {
	return repo.DB.Save(log).Error
}

func (repo PointLimitLogRepository) Create(log *business.PointLimitLog) error {
	return repo.DB.Create(log).Error
}
func (repo PointLimitLogRepository) FindLimitLog(by FindPointLimitLogBy) business.PointLimitLog {
	log := business.PointLimitLog{}
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
