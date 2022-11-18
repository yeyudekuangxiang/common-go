package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultDuiBaOrderRepository = DuiBaOrderRepository{DB: app.DB}

type DuiBaOrderRepository struct {
	DB *gorm.DB
}

func (repo DuiBaOrderRepository) Create(order *entity.DuiBaOrder) error {
	return repo.DB.Create(order).Error
}
func (repo DuiBaOrderRepository) Save(order *entity.DuiBaOrder) error {
	return repo.DB.Save(order).Error
}
func (repo DuiBaOrderRepository) FindByOrderId(orderId string) entity.DuiBaOrder {
	order := entity.DuiBaOrder{}
	err := repo.DB.Where("order_id = ?", orderId).First(&order).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return order
}

func (repo DuiBaOrderRepository) FindByUid(userId int64) entity.DuiBaOrder {
	order := entity.DuiBaOrder{}
	err := repo.DB.Where("user_id = ?", userId).First(&order).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return order
}
