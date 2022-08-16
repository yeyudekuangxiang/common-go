package business

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity/business"
)

type CarbonCreditsRepository struct {
	ctx *context.MioContext
}

func NewCarbonCreditsRepository(ctx *context.MioContext) *CarbonCreditsRepository {
	return &CarbonCreditsRepository{ctx: ctx}
}

func (repo CarbonCreditsRepository) Create(credit *business.CarbonCredits) error {
	return repo.ctx.DB.Create(credit).Error
}
func (repo CarbonCreditsRepository) Save(credit *business.CarbonCredits) error {
	return repo.ctx.DB.Save(credit).Error
}
func (repo CarbonCreditsRepository) FindCredits(userId int64) business.CarbonCredits {
	credit := business.CarbonCredits{}
	err := repo.ctx.DB.Model(credit).
		Where("b_user_id = ?", userId).
		Take(&credit).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return credit
}
