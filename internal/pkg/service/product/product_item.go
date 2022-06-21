package product

import (
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/util"
)

var DefaultProductItemService = NewProductItemService(repository.DefaultProductItemRepository)

func NewProductItemService(repo repository.ProductItemRepository) ProductItemService {
	return ProductItemService{repo: repo}
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

func (srv ProductItemService) GetListBy(param GetProductItemListParam) []entity.ProductItem {
	return srv.repo.GetListBy(repository.GetProductItemListBy{
		ItemIds: param.ItemIds,
	})
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
func (srv ProductItemService) FindProductByItemId(itemId string) (*entity.ProductItem, error) {
	item := srv.repo.FindByItemId(itemId)
	return &item, nil
}
func (srv ProductItemService) ListToMap(list []entity.ProductItem) map[string]entity.ProductItem {
	productMap := make(map[string]entity.ProductItem)
	for _, productItem := range list {
		productMap[productItem.ProductItemId] = productItem
	}
	return productMap
}
