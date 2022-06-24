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

	sql := `(SELECT buser.id user_id,COALESCE(sum(log.value),0) "value" 
FROM (SELECT * 
		FROM "business_carbon_credits_log" 
		WHERE created_at >= ? 
		AND created_at <= ?) as log 
right join business_user as "buser" 
on log.b_user_id = buser.id 
WHERE buser.b_company_id = ? 
GROUP BY buser.id) as t
`
	db := repo.DB.Table(sql, by.StartTime, by.EndTime, by.CompanyId)

	list := make([]business.UserCarbonRank, 0)
	var total int64

	err := db.Count(&total).Limit(by.Limit).Offset(by.Offset).Order("value desc").Find(&list).Error
	return list, total, err
}

// GetActualDepartmentCarbonRank 获取部门排行榜
func (repo CarbonCreditsLogRepository) GetActualDepartmentCarbonRank(by GetActualDepartmentCarbonRankBy) ([]business.DepartCarbonRank, int64, error) {
	sql := `
(SELECT
	r.department_id,sum(r.value) as  "value"
FROM
	(
	SELECT
		depart.top_id department_id,
		COALESCE ( SUM ( log.VALUE ), 0 ) "value" 
	FROM
		business_user AS buser
		INNER JOIN business_department AS depart ON buser.b_department_id = depart.ID 
		LEFT JOIN ( SELECT * FROM business_carbon_credits_log WHERE created_at >= ? AND created_at <= ? ) AS log ON buser.ID = log.b_user_id 
	WHERE
		depart.top_id <> 0 
		AND buser.b_company_id = ? 
	GROUP BY
		depart.top_id UNION ALL
	SELECT
		depart.ID AS department_id,
		COALESCE ( SUM ( log.VALUE ), 0 ) "value" 
	FROM
		business_user AS buser
		INNER JOIN business_department AS depart ON buser.b_department_id = depart.ID 
		LEFT JOIN ( SELECT * FROM business_carbon_credits_log WHERE created_at >= ? AND created_at <= ? ) AS log ON buser.ID = log.b_user_id 
	WHERE
		depart.top_id = 0 
		AND buser.b_company_id = ?
	GROUP BY
	depart.ID 
	
) r 
GROUP BY
	r.department_id
) t
`
	db := repo.DB.Table(sql,
		by.StartTime, by.EndTime, by.CompanyId,
		by.StartTime, by.EndTime, by.CompanyId,
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
