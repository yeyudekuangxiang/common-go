package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultStepHistoryRepository = StepHistoryRepository{DB: app.DB}

type StepHistoryRepository struct {
	DB *gorm.DB
}

func (repo StepHistoryRepository) FindBy(by FindStepHistoryBy) entity.StepHistory {
	sh := entity.StepHistory{}
	db := repo.DB.Model(sh)

	if by.UserId != 0 {
		db.Where("user_id = ?", by.UserId)
	}
	if !by.Day.IsZero() {
		db.Where("recorded_time = ?", by.Day)
	}
	if by.RecordedEpoch != 0 {
		db.Where("recorded_epoch = ?", by.RecordedEpoch)
	}

	repo.orderBy(db, by.OrderBy)

	err := db.First(&sh).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	return sh
}
func (repo StepHistoryRepository) orderBy(db *gorm.DB, orderByList entity.OrderByList) {
	for _, orderBy := range orderByList {
		switch orderBy {
		case entity.OrderByStepHistoryTimeDesc:
			db.Order("recorded_time desc")
		case entity.OrderByStepHistoryCountDesc:
			db.Order("count desc")
		}
	}
	return
}
func (repo StepHistoryRepository) GetStepHistoryList(by GetStepHistoryListBy) []entity.StepHistory {
	list := make([]entity.StepHistory, 0)

	db := repo.initStepHistoryListDB(repo.DB, by)

	err := db.Find(&list).Error
	if err != nil {
		panic(err)
	}
	return list
}
func (repo StepHistoryRepository) initStepHistoryListDB(db *gorm.DB, by GetStepHistoryListBy) *gorm.DB {
	db = db.Model(entity.StepHistory{})
	if by.UserId != 0 {
		db.Where("user_id = ?", by.UserId)
	}
	if len(by.RecordedEpochs) != 0 {
		db.Where("recorded_epoch in (?)", by.RecordedEpochs)
	}
	if !by.StartRecordedTime.IsZero() {
		db.Where("recorded_time >= ?", by.StartRecordedTime)
	}
	if !by.EndRecordedTime.IsZero() {
		db.Where("recorded_time <= ?", by.EndRecordedTime)
	}

	repo.orderBy(db, by.OrderBy)

	return db
}
func (repo StepHistoryRepository) GetStepHistoryPageList(by GetStepHistoryPageListBy) ([]entity.StepHistory, int64) {
	list := make([]entity.StepHistory, 0)

	db := repo.initStepHistoryListDB(repo.DB, by.GetStepHistoryListBy)
	var total int64
	err := db.Count(&total).Limit(by.Limit).Offset(by.Offset).Find(&list).Error
	if err != nil {
		panic(err)
	}
	return list, total
}
func (repo StepHistoryRepository) Save(history *entity.StepHistory) error {
	return repo.DB.Save(history).Error
}
func (repo StepHistoryRepository) Create(history *entity.StepHistory) error {
	return repo.DB.Create(history).Error
}

// GetUserLifeStepInfo 获取用户历史总步数及总天数
func (repo StepHistoryRepository) GetUserLifeStepInfo(userId int64) (steps int64, days int64) {
	sum := struct {
		Sum int64
	}{}
	err := repo.DB.Model(entity.StepHistory{}).Where("user_id = ?", userId).Count(&days).Select("sum(count) as sum").Take(&sum).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return sum.Sum, days
}
