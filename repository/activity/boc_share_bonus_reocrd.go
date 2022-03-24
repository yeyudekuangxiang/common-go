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
func (b BocShareBonusRecordRepository) GetUserBonus(userId int64, t activity.BocShareBonusType) int64 {
	type Sum struct {
		Sum int64
	}
	sum := Sum{}
	err := b.DB.Raw("select sum(value) as sum from boc_share_bonus_record where user_id = ? and type = ?", userId, t).
		First(&sum).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return sum.Sum
}
