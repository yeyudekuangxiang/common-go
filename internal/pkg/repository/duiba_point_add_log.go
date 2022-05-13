package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultDuiBaPointAddLogRepository = DuiBaPointAddLogRepository{DB: app.DB}

type DuiBaPointAddLogRepository struct {
	DB *gorm.DB
}

func (repo DuiBaPointAddLogRepository) Create(log *entity.DuiBaPointAddLog) error {
	return repo.DB.Create(log).Error
}
func (repo DuiBaPointAddLogRepository) Save(log *entity.DuiBaPointAddLog) error {
	return repo.DB.Save(log).Error
}
func (repo DuiBaPointAddLogRepository) FindBy(by FindDuiBaPointAddLogBy) entity.DuiBaPointAddLog {
	log := entity.DuiBaPointAddLog{}
	db := repo.DB.Model(log)

	if by.OrderNum != "" {
		db.Where("order_num = ?", by.OrderNum)
	}

	err := db.Take(&log).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	return log
}
func (repo DuiBaPointAddLogRepository) FindByID(id int64) entity.DuiBaPointAddLog {
	log := entity.DuiBaPointAddLog{}
	err := repo.DB.Take(&log, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	return log
}
