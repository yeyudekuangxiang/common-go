package activity

import (
	"gorm.io/gorm"
	"mio/core/app"
	"mio/model/entity/activity"
)

var DefaultBocShareBonusRecordRepository = BocShareBonusRecordRepository{DB: app.DB}

type BocShareBonusRecordRepository struct {
	DB *gorm.DB
}

func (b BocShareBonusRecordRepository) Save(record *activity.BocShareBonusRecord) error {
	return b.DB.Save(record).Error
}
