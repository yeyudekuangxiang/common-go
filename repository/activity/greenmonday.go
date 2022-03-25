package activity

import (
	"gorm.io/gorm"
	"mio/core/app"
	"mio/model/entity/activity"
)

var DefaultGMRecordRepository = GMRecordRepository{DB: app.DB}
var DefaultGMQuestionLogRepository = GMQuestionLogRepository{DB: app.DB}
var DefaultGMInvitationRecordRepository = GMInvitationRecordRepository{DB: app.DB}

type GMRecordRepository struct {
	DB *gorm.DB
}

func (repo GMRecordRepository) Save(record *activity.GMRecord) error {
	return repo.DB.Save(record).Error
}
func (repo GMRecordRepository) Create(record *activity.GMRecord) error {
	return repo.DB.Create(record).Error
}
func (repo GMRecordRepository) FindById(id int) activity.GMRecord {
	record := activity.GMRecord{}
	if err := repo.DB.First(&record, id).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return record
}
func (repo GMRecordRepository) FindBy(by FindGMRecordBy) activity.GMRecord {
	record := activity.GMRecord{}
	db := app.DB.Model(activity.GMRecord{})
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

type GMQuestionLogRepository struct {
	DB *gorm.DB
}

func (repo GMQuestionLogRepository) Save(log *activity.GMQuestionLog) error {
	return repo.DB.Save(log).Error
}
func (repo GMQuestionLogRepository) Create(log *activity.GMQuestionLog) error {
	return repo.DB.Create(log).Error
}
func (repo GMQuestionLogRepository) FindById(id int) activity.GMQuestionLog {
	log := activity.GMQuestionLog{}
	if err := repo.DB.First(&log, id).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return log
}
func (repo GMQuestionLogRepository) FindBy(by FindGMQuesLogBy) activity.GMQuestionLog {
	log := activity.GMQuestionLog{}
	db := app.DB.Model(activity.GMQuestionLog{})
	if by.UserId > 0 {
		db.Where("user_id = ?", by.UserId)
	}
	if err := db.First(&log).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return log
}

type GMInvitationRecordRepository struct {
	DB *gorm.DB
}

func (repo GMInvitationRecordRepository) Save(record *activity.GMInvitationRecord) error {
	return repo.DB.Save(record).Error
}
func (repo GMInvitationRecordRepository) Create(record *activity.GMInvitationRecord) error {
	return repo.DB.Create(record).Error
}
func (repo GMInvitationRecordRepository) FindById(id int) activity.GMInvitationRecord {
	record := activity.GMInvitationRecord{}
	if err := repo.DB.First(&record, id).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return record
}
func (repo GMInvitationRecordRepository) FindBy(by FindGMInvitationRecordBy) activity.GMInvitationRecord {
	record := activity.GMInvitationRecord{}
	db := app.DB.Model(activity.GMInvitationRecord{})
	if by.UserId > 0 {
		db.Where("user_id = ?", by.UserId)
	}
	if by.InviteeUserId > 0 {
		db.Where("invitee_user_id = ?", by.InviteeUserId)
	}
	if err := db.First(&record).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return record
}
