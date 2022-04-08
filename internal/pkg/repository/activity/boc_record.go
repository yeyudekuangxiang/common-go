package activity

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	activityM "mio/internal/pkg/model/entity/activity"
)

var DefaultBocRecordRepository = BocRecordRepository{
	DB: app.DB,
}

type BocRecordRepository struct {
	DB *gorm.DB
}

func (a BocRecordRepository) Save(record *activityM.BocRecord) error {
	return a.DB.Save(record).Error
}
func (a BocRecordRepository) FindBy(by FindRecordBy) activityM.BocRecord {
	record := activityM.BocRecord{}
	db := a.DB.Model(record)
	if by.UserId > 0 {
		db.Where("user_id = ?", by.UserId)
	}
	if err := db.First(&record).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return record
}
func (a BocRecordRepository) GetPageListBy(by GetRecordListBy) (list []activityM.BocRecord, total int64) {
	list = make([]activityM.BocRecord, 0)
	db := a.DB.Model(activityM.BocRecord{})
	if len(by.UserIds) > 0 {
		db.Where("user_id in (?)", by.UserIds)
	}

	if len(by.ShareUserIds) > 0 {
		db.Where("share_user_id in (?)", by.ShareUserIds)
	}

	if by.ApplyStatus > 0 {
		db.Where("apply_status = ?", by.ApplyStatus)
	}
	if by.ShareUserBocBonusStatus > 0 {
		db.Where("share_user_boc_bonus_status = ?", by.ShareUserBocBonusStatus)
	}

	db.Count(&total).Order("created_at desc").Offset(by.Offset).Limit(by.Limit)

	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return
}
