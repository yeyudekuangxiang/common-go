package business

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/business"
)

var DefaultCarbonCreditsRepository = CarbonCreditsRepository{DB: app.DB}

type CarbonCreditsRepository struct {
	DB *gorm.DB
}

func (repo CarbonCreditsRepository) Create(credit *business.CarbonCredits) error {
	return repo.DB.Create(credit).Error
}
func (repo CarbonCreditsRepository) Save(credit *business.CarbonCredits) error {
	return repo.DB.Create(credit).Error
}
func (repo CarbonCreditsRepository) FindCredits(userId int64) business.CarbonCredits {
	credit := business.CarbonCredits{}
	err := repo.DB.Model(credit).
		Where("b_user_id = ?", userId).
		Take(&credit).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return credit
}
