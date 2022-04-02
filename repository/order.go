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

// SubmitOrder 提交订单
func (repo OrderRepository) SubmitOrder(order *entity.Order, orderItems *[]entity.OrderItem) error {
	return repo.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		if err := tx.Create(orderItems).Error; err != nil {
			return err
		}
		return nil
	})
}
