package repository

import (
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/product"
)

var DefaultProductRepository IProductRepository = NewProductRepository()

type IProductRepository interface {
	ProductList() []Product.Product
}

func NewProductRepository() ProductRepository {
	return ProductRepository{}
}

type ProductRepository struct {
}

func (p ProductRepository) ProductList() []Product.Product {
	var products []Product.Product
	if err := app.DB.Table("product_item").Where("active = ?", "true").Order("sort desc,remaining_count desc").Find(&products).Error; err != nil {
		panic(err)
	}
	return products
}
