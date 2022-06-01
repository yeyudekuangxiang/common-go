package business

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/business"
)

var DefaultCarbonCreditsLogRepository = CarbonCreditsLogRepository{DB: app.DB}

type CarbonCreditsLogRepository struct {
	DB *gorm.DB
}

func (p CarbonCreditsLogRepository) Save(log *business.CarbonCreditsLog) error {
	return p.DB.Save(log).Error
}

func (p CarbonCreditsLogRepository) GetListBy(by GetCarbonCreditsLogListBy) []business.CarbonCreditsLog {
	list := make([]business.CarbonCreditsLog, 0)

	db := p.DB.Model(business.CarbonCreditsLog{})
	if by.UserId != 0 {
		db.Where("b_user_id = ?", by.UserId)
	}

	if !by.StartTime.IsZero() {
		db.Where("created_at >= ?", by.StartTime)
	}
	if !by.EndTime.IsZero() {
		db.Where("updated_at <= ?", by.EndTime)
	}

	if by.Type != "" {
		db.Where("type = ?", by.Type)
	}

	for _, orderBy := range by.OrderBy {
		switch orderBy {
		case business.OrderByCarbonCreditsLogCtDesc:
			db.Order("created_at desc")
		}
	}

	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}

	return list
}
