package business

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/business"
)

var DefaultPointRepository = PointRepository{DB: app.DB}

type PointRepository struct {
	DB *gorm.DB
}

func (repo PointRepository) Save(point *business.Point) error {
	return repo.DB.Save(point).Error
}
func (repo PointRepository) Create(point *business.Point) error {
	return repo.DB.Create(point).Error
}
func (repo PointRepository) FindPoint(by FindPointBy) business.Point {
	point := business.Point{}
	db := repo.DB.Model(point)

	if by.UserId != 0 {
		db.Where("b_user_id = ?", by.UserId)
	}

	err := db.Take(&point).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return point
}
