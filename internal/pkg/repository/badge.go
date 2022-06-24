package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository/repo_types"
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
func (repo BadgeRepository) FindUserCertCount(openid string) (int64, error) {
	var total int64
	return total, repo.DB.Model(entity.Badge{}).Where("openid = ?", openid).Count(&total).Error
}
func (repo BadgeRepository) FindBadge(by repo_types.FindBadgeBy) (*entity.Badge, error) {
	badge := entity.Badge{}
	db := repo.DB.Model(badge)
	if by.OrderId != "" {
		db.Where("order_id = ?", by.OrderId)
	}
	if by.ID != 0 {
		db.Where("id = ?", by.ID)
	}
	err := db.Take(&badge).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &badge, nil
}
