package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository/repo_types"
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
func (repo OrderRepository) FindByOrderId(orderId string) entity.Order {
	order := entity.Order{}
	err := repo.DB.Where("order_id = ?", orderId).First(&order).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return order
}
func (repo OrderRepository) GetPageFullOrder(do repo_types.GetPageFullOrderDO) ([]entity.OrderWithGood, int64, error) {
	db := repo.DB.Model(entity.OrderWithGood{})
	if do.Openid != "" {
		db.Where("openid = ?", do.Openid)
	}
	if do.OrderSource != "" {
		db.Where("source = ?", do.OrderSource)
	}
	list := make([]entity.OrderWithGood, 0)
	var total int64

	return list, total, db.Count(&total).Offset(do.Offset).Limit(do.Limit).Preload("OrderGoods").Find(&list).Error
}
