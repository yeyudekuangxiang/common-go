package repository

import (
	"errors"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultProductItemRepository = NewProductItemRepository(app.DB)

func NewProductItemRepository(db *gorm.DB) ProductItemRepository {
	return ProductItemRepository{DB: db}
}

type ProductItemRepository struct {
	DB *gorm.DB
}

func (repo ProductItemRepository) Save(product *entity.ProductItem) error {
	return repo.DB.Save(product).Error
}
func (repo ProductItemRepository) Create(product *entity.ProductItem) error {
	return repo.DB.Create(product).Error
}
func (repo ProductItemRepository) FindById(id int) entity.ProductItem {
	product := entity.ProductItem{}
	err := repo.DB.Where("id = ?", id).First(&product).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return product
}
func (repo ProductItemRepository) GetListBy(by GetProductItemListBy) []entity.ProductItem {
	list := make([]entity.ProductItem, 0)

	db := repo.DB.Model(entity.ProductItem{})
	if len(by.ItemIds) > 0 {
		db.Where("product_item_id in (?)", by.ItemIds)
	}
	err := db.Find(&list).Error
	if err != nil {
		panic(err)
	}
	return list
}

// CheckAndLockStock 检查并且锁定库存
func (repo ProductItemRepository) CheckAndLockStock(items []CheckStockItem) error {
	itemIds := make([]string, 0)
	itemMap := make(map[string]CheckStockItem)
	for _, item := range items {
		itemIds = append(itemIds, item.ItemId)
		itemMap[item.ItemId] = item
	}

	list := repo.GetListBy(GetProductItemListBy{
		ItemIds: itemIds,
	})
	if len(list) != len(itemIds) {
		return errors.New("存在失效商品,请去掉失效商品后重试")
	}
	return repo.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range list {
			wantCount := itemMap[item.ProductItemId].Count
			if !item.Active {
				return errors.New("商品`" + item.Title + "`已下架")
			}
			if item.RemainingCount < wantCount {
				return errors.New("商品`" + item.Title + "`库存不足")
			}
		}
		for _, item := range list {
			wantCount := itemMap[item.ProductItemId].Count
			item.RemainingCount -= wantCount
			err := tx.Save(&item).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// UnLockStock 释放库存
func (repo ProductItemRepository) UnLockStock(items []CheckStockItem) error {
	itemIds := make([]string, 0)
	itemMap := make(map[string]CheckStockItem)
	for _, item := range items {
		itemIds = append(itemIds, item.ItemId)
		itemMap[item.ItemId] = item
	}

	list := repo.GetListBy(GetProductItemListBy{
		ItemIds: itemIds,
	})

	return repo.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range list {
			unlockCount := itemMap[item.ProductItemId].Count
			item.RemainingCount += unlockCount
			err := tx.Save(&item).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}
