package repository

import (
	"gorm.io/gorm"
	"mio/core/app"
	"mio/model/entity"
)

var DefaultOrderRepository = NewOrderRepository(app.DB)

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return OrderRepository{DB: db}
}

type OrderRepository struct {
	DB *gorm.DB
}

func (repo OrderRepository) Save(order *entity.Order) error {
	return repo.DB.Save(order).Error
}
func (repo OrderRepository) Create(order *entity.Order) error {
	return repo.DB.Create(order).Error
}
