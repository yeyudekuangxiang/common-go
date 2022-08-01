package business

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/business"
)

var DefaultPointLogRepository = NewPointLogRepository(app.DB)

func NewPointLogRepository(db *gorm.DB) PointLogRepository {
	return PointLogRepository{
		DB: db,
	}
}

type PointLogRepository struct {
	DB *gorm.DB
}

func (p PointLogRepository) Save(log *business.PointLog) error {
	return p.DB.Save(log).Error
}
func (p PointLogRepository) Create(log *business.PointLog) error {
	return p.DB.Create(log).Error
}
func (p PointLogRepository) GetListBy(by GetPointLogListBy) []business.PointLog {
	list := make([]business.PointLog, 0)

	db := p.DB.Model(business.PointLog{})
	if by.UserId != 0 {
		db.Where("b_user_id = ?", by.UserId)
	}

	if !by.StartTime.IsZero() {
		db.Where("created_at >= ?", by.StartTime)
	}
	if !by.EndTime.IsZero() {
		db.Where("created_at <= ?", by.EndTime)
	}

	if by.Type != "" {
		db.Where("type = ?", by.Type)
	}

	for _, orderBy := range by.OrderBy {
		switch orderBy {
		case business.OrderByPointLogCTDESC:
			db.Order("created_at desc")
		}
	}

	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}

	return list
}

func (p PointLogRepository) GetUserTotalPoints(by GetCarbonCreditsLogSortedListBy) GetUserTotalCarbonCredits {
	list := GetUserTotalCarbonCredits{}
	db := p.DB.Model(business.PointLog{})
	if by.UserId != 0 {
		db.Where("b_user_id = ?", by.UserId)
	}
	if err := db.Select("sum(\"value\") as total").Find(&list).Error; err != nil {
		panic(err)
	}
	return list
}

//func (p PointTransactionRepository) GetSortListBy(by GetPointTransactionListBy) []business.PointTransaction {
//	list := make([]business.PointTransaction, 0)
//
//	db := p.DB.Model(business.PointTransaction{})
//	if by.UserId != 0 {
//		db.Where("b_user_id = ?", by.UserId)
//	}
//
//	if !by.StartTime.IsZero() {
//		db.Where("created_at >= ?", by.StartTime)
//	}
//	if !by.EndTime.IsZero() {
//		db.Where("updated_at <= ?", by.EndTime)
//	}
//
//	if by.Type != "" {
//		db.Where("type = ?", by.Type)
//	}
//
//	for _, orderBy := range by.OrderBy {
//		switch orderBy {
//		case business.OrderByPointTranCTDESC:
//			db.Order("created_at desc")
//		}
//	}
//
//	if err := db.Select("sum(value) as total ,type ").Group("type").Order("total desc").Find(&list).Error; err != nil {
//		panic(err)
//	}
//
//	return list
//}
