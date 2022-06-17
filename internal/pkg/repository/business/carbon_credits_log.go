package business

import (
	"fmt"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/business"
)

var DefaultCarbonCreditsLogRepository = CarbonCreditsLogRepository{DB: app.DB}

type CarbonCreditsLogRepository struct {
	DB *gorm.DB
}

func (repo CarbonCreditsLogRepository) Save(log *business.CarbonCreditsLog) error {
	return repo.DB.Save(log).Error
}
func (repo CarbonCreditsLogRepository) Create(log *business.CarbonCreditsLog) error {
	return repo.DB.Create(log).Error
}
func (repo CarbonCreditsLogRepository) GetListBy(by GetCarbonCreditsLogListBy) []business.CarbonCreditsLog {
	list := make([]business.CarbonCreditsLog, 0)

	db := repo.DB.Model(business.CarbonCreditsLog{})
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

// GetActualUserCarbonRank 获取用户碳积分排行榜
func (repo CarbonCreditsLogRepository) GetActualUserCarbonRank(by GetActualUserCarbonRankBy) ([]business.UserCarbonRank, int64, error) {
	db := repo.DB.Table(fmt.Sprintf("%s as log", business.CarbonCreditsLog{}.TableName())).
		Joins(fmt.Sprintf("inner join %s as \"buser\" on log.b_user_id = buser.id", business.User{}.TableName()))

	if by.CompanyId != 0 {
		db.Where("buser.b_company_id = ?", by.CompanyId)
	}
	if !by.StartTime.IsZero() {
		db.Where("log.created_at >= ?", by.StartTime)
	}
	if !by.EndTime.IsZero() {
		db.Where("log.created_at <= ?", by.EndTime)
	}

	db.Select("log.b_user_id user_id,sum(log.value) \"value\"").Group("log.b_user_id")

	db = repo.DB.Table("(?) t", db)

	list := make([]business.UserCarbonRank, 0)
	var total int64

	err := db.Count(&total).Limit(by.Limit).Offset(by.Offset).Order("value desc").Find(&list).Error
	return list, total, err
}

// GetActualDepartmentCarbonRank 获取部门排行榜
func (repo CarbonCreditsLogRepository) GetActualDepartmentCarbonRank(by GetActualDepartmentCarbonRankBy) ([]business.DepartCarbonRank, int64, error) {
	db := repo.DB.Table("(SELECT depart.top_id  department_id,sum(log.value) \"value\" "+
		"FROM business_carbon_credits_log AS log "+
		"INNER JOIN business_user AS buser ON log.b_user_id = buser.ID "+
		"INNER JOIN business_department as depart on buser.b_department_id = depart.id "+
		"where depart.top_id <> 0 and buser.b_company_id = ? and log.created_at >= ? and log.created_at <= ? "+
		"GROUP BY depart.top_id "+
		"UNION SELECT depart.id as department_id,sum(log.value) \"value\" "+
		"FROM business_carbon_credits_log AS log I"+
		"NNER JOIN business_user AS buser ON log.b_user_id = buser.ID "+
		"INNER JOIN business_department as depart on buser.b_department_id = depart.id  "+
		"where depart.top_id = 0 and buser.b_company_id = ? and log.created_at >= ? and log.created_at <= ? "+
		"GROUP BY depart.id) t",
		by.CompanyId, by.StartTime, by.EndTime,
		by.CompanyId, by.StartTime, by.EndTime,
	)

	list := make([]business.DepartCarbonRank, 0)
	var total int64

	err := db.Count(&total).Limit(by.Limit).Offset(by.Offset).Order("value desc").Find(&list).Error
	return list, total, err
}
func (repo CarbonCreditsLogRepository) GetSortedListBy(by GetCarbonCreditsLogSortedListBy) []CarbonCreditsLogSortedList {
	list := make([]CarbonCreditsLogSortedList, 0)

	db := repo.DB.Model(business.CarbonCreditsLog{})
	if by.UserId != 0 {
		db.Where("b_user_id = ?", by.UserId)
	}
	if !by.StartTime.IsZero() {
		db.Where("created_at >= ?", by.StartTime)
	}
	//纠正了已出错误 updated_at
	if !by.EndTime.IsZero() {
		db.Where("created_at <= ?", by.EndTime)
	}
	if len(by.UserIds) > 0 {
		db.Where("b_user_id in (?)", by.UserIds)
	}

	if err := db.Select("sum(value) as total ,type ").Group("type").Order("total desc").Find(&list).Error; err != nil {
		panic(err)
	}

	return list
}

func (repo CarbonCreditsLogRepository) GetCarbonCreditsLogListHistory(by GetCarbonCreditsLogSortedListBy) []CarbonCreditsLogListHistory {
	list := make([]CarbonCreditsLogListHistory, 0)
	db := repo.DB.Model(business.CarbonCreditsLog{})
	if by.UserId != 0 {
		db.Where("b_user_id = ?", by.UserId)
	}
	if err := db.Select("sum(\"value\") as total ,substring(cast(created_at as varchar),1,7) as month,type").Group("month,type").Order("month desc,total desc").Find(&list).Error; err != nil {
		panic(err)
	}
	return list
}

func (repo CarbonCreditsLogRepository) GetUserTotalCarbonCredits(by GetCarbonCreditsLogSortedListBy) GetUserTotalCarbonCredits {
	list := GetUserTotalCarbonCredits{}
	db := repo.DB.Model(business.CarbonCreditsLog{})
	if by.UserId != 0 {
		db.Where("b_user_id = ?", by.UserId)
	}
	if err := db.Select("sum(\"value\") as total").Find(&list).Error; err != nil {
		panic(err)
	}
	return list
}
