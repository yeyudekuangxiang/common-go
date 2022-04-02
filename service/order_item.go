package service

import (
	"mio/model/entity"
	"mio/repository"
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
