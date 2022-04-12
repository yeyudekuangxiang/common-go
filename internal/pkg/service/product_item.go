package service

import (
	"mio/internal/pkg/model/entity"
	repository2 "mio/internal/pkg/repository"
)

var DefaultProductItemService = NewProductItemService(repository2.DefaultProductItemRepository)

func NewProductItemService(repo repository2.ProductItemRepository) ProductItemService {
	return ProductItemService{repo: repo}
}

type ProductItemService struct {
	repo repository2.ProductItemRepository
}

// CheckAndLockStock 检查并锁定库存
func (srv ProductItemService) CheckAndLockStock(items []repository2.CheckStockItem) error {
	return srv.repo.CheckAndLockStock(items)
}

// UnLockStock 释放下单失败的库存
func (srv ProductItemService) UnLockStock(items []repository2.CheckStockItem) error {
	return srv.repo.UnLockStock(items)
}

func (srv ProductItemService) GetListBy(by repository2.GetProductItemListBy) []entity.ProductItem {
	return srv.repo.GetListBy(by)
}
