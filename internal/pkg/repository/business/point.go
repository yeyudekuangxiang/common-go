package business

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity/business"
)

type PointRepository struct {
	ctx *context.MioContext
}

func NewPointRepository(ctx *context.MioContext) *PointRepository {
	return &PointRepository{ctx: ctx}
}

func (repo PointRepository) Save(point *business.Point) error {
	return repo.ctx.DB.Save(point).Error
}
func (repo PointRepository) Create(point *business.Point) error {
	return repo.ctx.DB.Create(point).Error
}
func (repo PointRepository) FindPoint(by FindPointBy) business.Point {
	point := business.Point{}
	db := repo.ctx.DB.Model(point)

	if by.UserId != 0 {
		db.Where("b_user_id = ?", by.UserId)
	}

	err := db.Take(&point).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return point
}
