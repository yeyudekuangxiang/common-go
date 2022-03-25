package service

import (
	"gorm.io/gorm"
	"mio/model/entity"
	"mio/repository"
)

var DefaultProductItemService = NewProductItemService(repository.DefaultProductItemRepository)

func NewProductItemService(repo repository.ProductItemRepository) ProductItemService {
	return ProductItemService{repo: repo}
}
func NewProductItemServiceByDB(db *gorm.DB) ProductItemService {
	return NewProductItemService(repository.NewProductItemRepository(db))
}

type ProductItemService struct {
	repo repository.ProductItemRepository
}

// CheckAndLockStock 检查并锁定库存
func (srv ProductItemService) CheckAndLockStock(items []repository.CheckStockItem) error {
	return srv.repo.CheckAndLockStock(items)
}

// UnLockStock 释放下单失败的库存
func (srv ProductItemService) UnLockStock(items []repository.CheckStockItem) error {
	return srv.repo.UnLockStock(items)
}

func (srv ProductItemService) GetListBy(by repository.GetProductItemListBy) []entity.ProductItem {
	return srv.repo.GetListBy(by)
}
