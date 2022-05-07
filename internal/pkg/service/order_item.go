package service

import (
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

var DefaultOrderItemService = NewOrderItemService(repository.DefaultOrderItemRepository)

func NewOrderItemService(repo repository.OrderItemRepository) OrderItemService {
	return OrderItemService{repo: repo}
}

type OrderItemService struct {
	repo repository.OrderItemRepository
}

func (srv OrderItemService) CreateOrderItems(orderId string, submitOrderItems []submitOrderItem) error {
	orderItems := make([]entity.OrderItem, 0)
	for _, submitOrderItem := range submitOrderItems {
		orderItems = append(orderItems, entity.OrderItem{
			OrderId: orderId,
			ItemId:  submitOrderItem.ItemId,
			Count:   submitOrderItem.Count,
			Cost:    submitOrderItem.Cost,
		})
	}
	return srv.repo.CreateBatch(&orderItems)
}
func (srv OrderItemService) GetOrderItemListByOrderId(orderId string) ([]entity.OrderItem, error) {
	return srv.repo.GetListByOrderId(orderId), nil
}
