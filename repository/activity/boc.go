package activity

import (
	"gorm.io/gorm"
	"mio/core/app"
	activityM "mio/model/entity/activity"
)

var DefaultBocApplyRecordRepository = BocApplyRecordRepository{
	DB: app.DB,
}

type BocApplyRecordRepository struct {
	DB *gorm.DB
}

func (a BocApplyRecordRepository) Save(record *activityM.BocApplyRecord) error {
	return a.DB.Save(record).Error
}
func (a BocApplyRecordRepository) FindBy(by FindRecordBy) activityM.BocApplyRecord {
	record := activityM.BocApplyRecord{}
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
func (a BocApplyRecordRepository) GetPageListBy(by GetRecordListBy) (list []activityM.BocApplyRecord, total int64) {
	list = make([]activityM.BocApplyRecord, 0)
	db := a.DB.Model(activityM.BocApplyRecord{})
	if len(by.UserIds) > 0 {
		db.Where("user_id in (?)", by.UserIds)
	}

	if len(by.ShareUserIds) > 0 {
		db.Where("share_user_id in (?)", by.ShareUserIds)
	}

	if by.ApplyStatus > 0 {
		db.Where("apply_status = ?", by.ApplyStatus)
	}

	db.Count(&total).Order("created_at desc").Offset(by.Offset).Limit(by.Limit)

	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return
}
