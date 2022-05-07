package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultOrderItemRepository = NewOrderItemRepository(app.DB)

func NewOrderItemRepository(db *gorm.DB) OrderItemRepository {
	return OrderItemRepository{DB: db}
}

type OrderItemRepository struct {
	DB *gorm.DB
}

func (repo OrderItemRepository) Save(item *entity.OrderItem) error {
	return repo.DB.Save(item).Error
}
func (repo OrderItemRepository) Create(item *entity.OrderItem) error {
	return repo.DB.Create(item).Error
}
func (repo OrderItemRepository) CreateBatch(item *[]entity.OrderItem) error {
	return repo.DB.Create(item).Error
}
func (repo OrderItemRepository) GetListByOrderId(orderId string) []entity.OrderItem {
	list := make([]entity.OrderItem, 0)
	err := repo.DB.Where("order_id = ?", orderId).Find(&list).Error
	if err != nil {
		panic(err)
	}
	return list
}
