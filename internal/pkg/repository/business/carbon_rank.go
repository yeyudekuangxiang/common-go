package business

import (
	"fmt"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/business"
	"mio/pkg/errno"
	"time"
)

var DefaultCarbonRankRepository = CarbonRankRepository{DB: app.DB}

type CarbonRankRepository struct {
	DB *gorm.DB
}

// CreateBatch 批量创建
func (repo CarbonRankRepository) CreateBatch(list *[]business.CarbonRank) error {
	return repo.DB.Create(list).Error
}

// FindCarbonRank 根据条件查询一条排名信息
func (repo CarbonRankRepository) FindCarbonRank(by FindCarbonRankBy) business.CarbonRank {
	rank := business.CarbonRank{}
	db := repo.DB.Model(rank)

	if by.Pid != 0 {
		db.Where("pid = ?", by.Pid)
	}
	if by.ObjectType != "" {
		db.Where("object_type = ?", by.ObjectType)
	}
	if by.DateType != "" {
		db.Where("date_type = ?", by.DateType)
	}
	fmt.Println("失去", time.Now())
	if !by.TimePoint.IsZero() {
		db.Where("time_point =  ?", by.TimePoint)
	}
	err := db.Take(&rank).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return rank
}

// Save 保存
func (repo CarbonRankRepository) Save(rank *business.CarbonRank) error {
	return repo.DB.Save(rank).Error
}

// Create 创建
func (repo CarbonRankRepository) Create(rank *business.CarbonRank) error {
	return repo.DB.Create(rank).Error
}

// GetCarbonUserRankList 获取碳积分排行榜
func (repo CarbonRankRepository) GetCarbonUserRankList(by GetCarbonRankBy) ([]business.CarbonRank, int64, error) {
	db := repo.DB.Table(fmt.Sprintf("%s as rank", business.CarbonRank{}.TableName())).
		Joins(fmt.Sprintf("inner join %s as \"buser\" on rank.pid = buser.id", business.User{}.TableName())).
		Where("rank.object_type = 'user'", by.ObjectType)

	if by.DateType != "" {
		db.Where("rank.date_type = ?", by.DateType)
	}
	if by.CompanyId != 0 {
		db.Where("buser.b_company_id = ?", by.CompanyId)
	}
	if !by.TimePoint.IsZero() {
		db.Where("rank.time_point = ?", by.TimePoint)
	}

	list := make([]business.CarbonRank, 0)
	var total int64

	err := db.Count(&total).
		Limit(by.Limit).
		Offset(by.Offset).
		Order("rank asc").
		Find(&list).Error
	return list, total, err
}

// GetCarbonRankList 获取碳积分排行榜
func (repo CarbonRankRepository) GetCarbonRankList(by GetCarbonRankBy) ([]business.CarbonRank, int64, error) {
	objectTable := ""
	if by.ObjectType == "user" {
		objectTable = business.User{}.TableName()
	} else if by.ObjectType == "department" {
		objectTable = business.Department{}.TableName()
	} else {
		return nil, 0, errno.ErrInternalServer
	}

	db := repo.DB.Table(fmt.Sprintf("%s as rank", business.CarbonRank{}.TableName())).
		Joins(fmt.Sprintf("inner join %s as \"object\" on rank.pid = object.id", objectTable)).
		Where("rank.object_type = ?", by.ObjectType)

	if by.DateType != "" {
		db.Where("rank.date_type = ?", by.DateType)
	}
	if by.CompanyId != 0 {
		db.Where("object.b_company_id = ?", by.CompanyId)
	}
	if !by.TimePoint.IsZero() {
		db.Where("rank.time_point = ?", by.TimePoint)
	}

	list := make([]business.CarbonRank, 0)
	var total int64

	err := db.Count(&total).
		Limit(by.Limit).
		Offset(by.Offset).
		Order("rank asc").
		Find(&list).Error
	return list, total, err
}
