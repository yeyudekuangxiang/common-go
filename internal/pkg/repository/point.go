package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
)

func NewPointRepository(ctx *context.MioContext) *PointRepository {
	return &PointRepository{ctx: ctx}
}

type PointRepository struct {
	ctx *context.MioContext
}

func (repo PointRepository) Save(point *entity.Point) error {
	app.DB.Session(&gorm.Session{})
	return repo.ctx.DB.Save(point).Error
}

func (repo PointRepository) FindBy(by FindPointBy) entity.Point {
	point := entity.Point{}
	db := repo.ctx.DB.Model(point)
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

func (repo PointRepository) FindListPoint(by FindListPoint) []entity.Point {
	list := make([]entity.Point, 0)
	point := entity.Point{}
	db := repo.ctx.DB.Model(point)
	if len(by.OpenIds) != 0 {
		db.Where("openid in(?)", by.OpenIds)
	}
	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return list
}
func (repo PointRepository) FindForUpdate(openId string) (entity.Point, error) {
	point := entity.Point{}
	err := repo.ctx.DB.
		Set("gorm:query_option", "for update").
		Where("openid = ?", openId).
		First(&point).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return point, err
	}
	return point, nil
}
func (repo PointRepository) Create(point *entity.Point) error {
	return repo.ctx.DB.Create(point).Error
}
