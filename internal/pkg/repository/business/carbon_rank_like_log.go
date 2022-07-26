package business

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/business"
)

var DefaultCarbonRankLikeLogRepository = CarbonRankLikeLogRepository{DB: app.DB}

type CarbonRankLikeLogRepository struct {
	DB *gorm.DB
}

func (repo CarbonRankLikeLogRepository) Create(likeNum *business.CarbonRankLikeLog) error {
	return repo.DB.Create(likeNum).Error
}
func (repo CarbonRankLikeLogRepository) Save(likeNum *business.CarbonRankLikeLog) error {
	return repo.DB.Save(likeNum).Error
}
func (repo CarbonRankLikeLogRepository) FindLikeLog(by FindCarbonRankLikeLogBy) business.CarbonRankLikeLog {
	likeLog := business.CarbonRankLikeLog{}
	db := repo.DB.Model(likeLog)

	if !by.TimePoint.IsZero() {
		db.Where("time_point = ?", by.TimePoint)
	}
	if by.UserId != 0 {
		db.Where("b_user_id = ?", by.UserId)
	}
	if by.Pid != 0 {
		db.Where("pid = ?", by.Pid)
	}
	if by.DateType != "" {
		db.Where("date_type = ?", by.DateType)
	}
	if by.ObjectType != "" {
		db.Where("object_type = ?", by.ObjectType)
	}
	err := db.Take(&likeLog).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return likeLog
}
func (repo CarbonRankLikeLogRepository) GetLikeLogList(by GetCarbonRankLikeLogListBy) []business.CarbonRankLikeLog {
	db := repo.DB.Model(business.CarbonRankLikeLog{})

	if !by.TimePoint.IsZero() {
		db.Where("time_point = ?", by.TimePoint)
	}
	if by.UserId != 0 {
		db.Where("b_user_id = ?", by.UserId)
	}
	if len(by.PIds) > 0 {
		db.Where("pid in (?)", by.PIds)
	}
	if by.DateType != "" {
		db.Where("date_type = ?", by.DateType)
	}
	if by.ObjectType != "" {
		db.Where("object_type = ?", by.ObjectType)
	}

	list := make([]business.CarbonRankLikeLog, 0)
	err := db.Find(&list).Error

	if err != nil {
		panic(err)
	}
	return list
}
