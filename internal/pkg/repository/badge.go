package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultBadgeRepository = BadgeRepository{DB: app.DB}

type BadgeRepository struct {
	DB *gorm.DB
}

func (repo BadgeRepository) FindLastWithType(t entity.CertType) entity.Badge {
	badge := entity.Badge{}
	err := repo.DB.Model(badge).
		Joins("inner join certificate on badge.certificate_id = certificate.certificate_id").
		Where("certificate.type = ?", t).
		Where("badge.code is not null").
		Order("code desc").Take(&badge).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return badge
}
func (repo BadgeRepository) Create(badge *entity.Badge) error {
	return repo.DB.Create(badge).Error
}
