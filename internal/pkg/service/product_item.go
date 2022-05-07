package service

import (
	"mio/internal/pkg/model/entity"
	repository2 "mio/internal/pkg/repository"
	"mio/internal/pkg/util"
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
func (srv ProductItemService) CreateOrUpdateProductItem(param CreateOrUpdateProductItemParam) (*entity.ProductItem, error) {
	item := srv.repo.FindByItemId(param.ItemId)
	if item.ID == 0 {
		item = entity.ProductItem{
			ProductItemId:          param.ItemId,
			Virtual:                param.Virtual,
			Title:                  param.Title,
			Cost:                   param.Cost,
			ImageUrl:               param.ImageUrl,
			SalesCount:             0,
			Active:                 true,
			ProductItemReferenceId: util.UUID(),
			Sort:                   param.Sort,
		}
		return &item, srv.repo.Create(&item)
	}
	item.Virtual = param.Virtual
	item.Title = param.Title
	item.Cost = param.Cost
	item.ImageUrl = param.ImageUrl
	item.Sort = param.Sort
	return &item, srv.repo.Save(&item)
}
