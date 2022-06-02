package business

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/business"
)

var DefaultCarbonRankLikeNumRepository = CarbonRankLikeNumRepository{DB: app.DB}

type CarbonRankLikeNumRepository struct {
	DB *gorm.DB
}

func (repo CarbonRankLikeNumRepository) Create(likeNum *business.CarbonRankLikeNum) error {
	return repo.DB.Create(likeNum).Error
}
func (repo CarbonRankLikeNumRepository) Save(likeNum *business.CarbonRankLikeNum) error {
	return repo.DB.Save(likeNum).Error
}

func (repo CarbonRankLikeNumRepository) GetLikeNumList(by GetCarbonRankLikeNumListBy) []business.CarbonRankLikeNum {
	db := repo.DB.Model(business.CarbonRankLikeNum{})

	if !by.TimePoint.IsZero() {
		db.Where("time_point = ?", by.TimePoint)
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
	list := make([]business.CarbonRankLikeNum, 0)
	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return list
}
func (repo CarbonRankLikeNumRepository) FindLikeNum(by FindCarbonRankLikeNumBy) business.CarbonRankLikeNum {
	likeNum := business.CarbonRankLikeNum{}
	db := repo.DB.Model(likeNum)

	if !by.TimePoint.IsZero() {
		db.Where("time_point = ?", by.TimePoint)
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

	err := db.Take(&likeNum).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return likeNum
}
