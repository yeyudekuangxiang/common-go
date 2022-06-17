package activity

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/activity"
	"strings"
)

var DefaultGDDonationBookRepository = GDDonationBookRepository{DB: app.DB}
var DefaultGDDbUserSchoolRepository = GDDbUserSchoolRepository{DB: app.DB}
var DefaultGDDbSchoolRankRepository = GDDbSchoolRankRepository{DB: app.DB}
var DefaultGDDbCityRepository = GDDbCityRepository{DB: app.DB}
var DefaultGDDbSchoolRepository = GDDbSchoolRepository{DB: app.DB}
var DefaultGDDbGradeRepository = GDDbGradeRepository{DB: app.DB}

type GDDonationBookRepository struct {
	DB *gorm.DB
}

func (repo GDDonationBookRepository) Save(record *activity.GDDonationBookRecord) error {
	return repo.DB.Save(record).Error
}
func (repo GDDonationBookRepository) Create(record *activity.GDDonationBookRecord) error {
	return repo.DB.Create(record).Error
}
func (repo GDDonationBookRepository) FindById(id int64) activity.GDDonationBookRecord {
	record := activity.GDDonationBookRecord{}
	if err := repo.DB.First(&record, id).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return record
}
func (repo GDDonationBookRepository) FindBy(by FindRecordBy) activity.GDDonationBookRecord {
	record := activity.GDDonationBookRecord{}
	db := app.DB.Model(activity.GDDonationBookRecord{})
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

// GetUserBy 获取用户活动及个人信息
func (repo GDDonationBookRepository) GetUserBy(by FindRecordBy) activity.GDDonationBookRecord {
	var record activity.GDDonationBookRecord
	db := app.DB.Model(activity.GDDonationBookRecord{})
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

// GetInvitedBy 获取用户活动及个人信息
func (repo GDDonationBookRepository) GetInvitedBy(by FindRecordBy) []activity.GDDonationBookRecord {
	record := make([]activity.GDDonationBookRecord, 0)
	db := app.DB.Model(activity.GDDonationBookRecord{})
	if by.UserId > 0 {
		db.Where("gd_donation_book.invite_id = ?", by.UserId)
	}
	if err := db.Find(&record).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return record
}

type GDDbSchoolRankRepository struct {
	DB *gorm.DB
}

func (repo GDDbSchoolRankRepository) Save(record *activity.GDDbSchoolRank) error {
	return repo.DB.Save(record).Error
}
func (repo GDDbSchoolRankRepository) Create(record *activity.GDDbSchoolRank) error {
	return repo.DB.Create(record).Error
}
func (repo GDDbSchoolRankRepository) FindById(id int64) activity.GDDbSchoolRank {
	record := activity.GDDbSchoolRank{}
	if err := repo.DB.First(&record, id).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return record
}
func (repo GDDbSchoolRankRepository) FindBy(by FindSchoolBy) activity.GDDbSchoolRank {
	record := activity.GDDbSchoolRank{}
	db := app.DB.Model(activity.GDDbSchoolRank{})
	if len(by.SchoolIds) > 0 {
		db.Where("school_id in ?", by.SchoolIds)
	}
	if by.SchoolName != "" {
		db.Where("school_name like ?", strings.Join([]string{by.SchoolName, "%"}, ""))
	}
	if err := db.First(&record).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return record
}

func (repo GDDbSchoolRankRepository) FindAllBy(by FindSchoolBy) []activity.GDDbSchoolRank {
	record := make([]activity.GDDbSchoolRank, 0)
	db := app.DB.Model(activity.GDDbSchoolRank{})
	if len(by.SchoolIds) > 0 {
		db.Where("school_id in ?", by.SchoolIds)
	}
	if by.SchoolName != "" {
		db.Where("school_name like ?", strings.Join([]string{by.SchoolName, "%"}, ""))
	}
	if err := db.Find(&record).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return record
}

func (repo GDDbSchoolRankRepository) GetRank() []activity.GDDbSchoolRank {
	record := make([]activity.GDDbSchoolRank, 0)
	db := app.DB.Model(activity.GDDbSchoolRank{}).Order("donate_number desc").Limit(20)
	if err := db.Scan(&record).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return record
}

type GDDbSchoolRepository struct {
	DB *gorm.DB
}

func (repo GDDbSchoolRepository) Save(record *activity.GDDbSchool) error {
	return repo.DB.Save(record).Error
}

func (repo GDDbSchoolRepository) Create(record *activity.GDDbSchool) error {
	return repo.DB.Create(record).Error
}

func (repo GDDbSchoolRepository) FindById(id int64) activity.GDDbSchool {
	log := activity.GDDbSchool{}
	if err := repo.DB.First(&log, id).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return log
}

func (repo GDDbSchoolRepository) FindBy(by FindSchoolBy) activity.GDDbSchool {
	record := activity.GDDbSchool{}
	db := app.DB.Model(activity.GDDbSchool{})
	if len(by.SchoolIds) < 1 {
		db.Where("id in ?", by.SchoolIds)
	}
	if by.SchoolName != "" {
		db.Where("school_name like ?", strings.Join([]string{by.SchoolName, "%"}, ""))
	}
	if by.GradeType > 0 {
		db.Where("type = ?", by.GradeType).Or("type = ?", 0)
	}
	if by.CityId != 0 {
		db.Where("city_id = ?", by.CityId)
	}
	if err := db.First(&record).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return record
}

func (repo GDDbSchoolRepository) FindAllBy(by FindSchoolBy) []activity.GDDbSchool {
	record := make([]activity.GDDbSchool, 0)
	db := app.DB.Model(activity.GDDbSchool{})
	if len(by.SchoolIds) < 1 {
		db.Where("id in ?", by.SchoolIds)
	}
	if by.SchoolName != "" {
		db.Where("school_name like ?", strings.Join([]string{by.SchoolName, "%"}, ""))
	}
	if by.GradeType > 0 {
		db.Where("type = ?", by.GradeType).Or("type = ?", 0)
	}
	if by.CityId != 0 {
		db.Where("city_id = ?", by.CityId)
	}
	if err := db.Find(&record).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return record
}

type GDDbUserSchoolRepository struct {
	DB *gorm.DB
}

func (repo GDDbUserSchoolRepository) Save(record *activity.GDDbUserSchool) error {
	return repo.DB.Save(record).Error
}
func (repo GDDbUserSchoolRepository) Create(record *activity.GDDbUserSchool) error {
	return repo.DB.Create(record).Error
}
func (repo GDDbUserSchoolRepository) FindById(id int64) activity.GDDbUserSchool {
	record := activity.GDDbUserSchool{}
	if err := repo.DB.First(&record, id).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return record
}

func (repo GDDbUserSchoolRepository) FindBy(by FindRecordBy) activity.GDDbUserSchool {
	record := activity.GDDbUserSchool{}
	db := app.DB.Model(activity.GDDbUserSchool{})
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

type GDDbCityRepository struct {
	DB *gorm.DB
}

func (repo GDDbCityRepository) FindAll() []activity.GDDbCity {
	record := make([]activity.GDDbCity, 0)
	if err := app.DB.Model(activity.GDDbCity{}).Find(&record).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return record
}

func (repo GDDbCityRepository) FindById(id int64) activity.GDDbCity {
	record := activity.GDDbCity{}
	if err := repo.DB.First(&record, id).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return record
}

type GDDbGradeRepository struct {
	DB *gorm.DB
}

func (repo GDDbGradeRepository) FindAll() []activity.GDDbGrade {
	record := make([]activity.GDDbGrade, 0)
	if err := app.DB.Model(activity.GDDbGrade{}).Find(&record).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return record
}

func (repo GDDbGradeRepository) FindById(id int64) activity.GDDbGrade {
	record := activity.GDDbGrade{}
	if err := repo.DB.First(&record, id).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return record
}
