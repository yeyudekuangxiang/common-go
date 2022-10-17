package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository/repotypes"
)

func NewCityRepository(ctx *context.MioContext) CityRepository {
	return CityRepository{
		ctx: ctx,
	}
}

type CityRepository struct {
	ctx *context.MioContext
}

func (repo CityRepository) Save(day *entity.City) error {
	return repo.ctx.DB.Save(day).Error
}

func (repo CityRepository) Create(day *entity.City) error {
	return repo.ctx.DB.Create(day).Error
}

func (repo CityRepository) GetList(by repotypes.GetCityListDO) ([]entity.City, error) {
	list := make([]entity.City, 0)
	db := repo.ctx.DB.Model(entity.City{})
	if len(by.CityCodeSlice) != 0 {
		db.Where("city_code in (?)", by.CityCodeSlice)
	}
	if by.CityCode != "" {
		db.Where("city_code", by.CityCode)
	}
	if by.Name != "" {
		db.Where("name = ?", by.Name)
	}
	if err := db.Find(&list).Error; err != nil {
		panic(err)
	}
	return list, nil
}

func (repo CityRepository) GetByCityCode(by repotypes.GetCityByCode) (entity.City, error) {
	info := entity.City{}
	db := repo.ctx.DB.Model(info)
	if by.CityCode != "" {
		db.Where("city_code", by.CityCode)
	}
	if err := db.First(&info).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return info, nil
}
