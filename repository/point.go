package repository

import (
	"gorm.io/gorm"
	"mio/core/app"
	"mio/model/entity"
)

var DefaultPointRepository = NewPointRepository(app.DB)

func NewPointRepository(db *gorm.DB) PointRepository {
	return PointRepository{DB: db}
}

type PointRepository struct {
	DB *gorm.DB
}

func (p PointRepository) Save(point *entity.Point) error {
	return p.DB.Save(point).Error
}

func (p PointRepository) FindBy(by FindPointBy) entity.Point {
	point := entity.Point{}
	db := p.DB.Model(point)
	if by.OpenId != "" {
		db.Where("openid = ?", by.OpenId)
	}
	if err := db.First(&point).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return point
}
