package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultStepRepository = StepRepository{DB: app.DB}

type StepRepository struct {
	DB *gorm.DB
}

func (repo StepRepository) FindBy(by FindStepBy) entity.Step {
	step := entity.Step{}
	db := repo.DB.Model(step)
	if by.UserId != 0 {
		db.Where("user_id = ?", by.UserId)
	}

	err := db.First(&step).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return step
}
func (repo StepRepository) Save(step *entity.Step) error {
	return repo.DB.Save(step).Error
}
func (repo StepRepository) Create(step *entity.Step) error {
	return repo.DB.Create(step).Error
}
