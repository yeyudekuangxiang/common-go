package product

import (
	"mio/internal/pkg/model/entity/product"
	"mio/internal/pkg/repository"
)

var DefaultProductService = NewProductService(repository.DefaultProductRepository)

func NewProductService(r repository.IProductRepository) ProductService {
	return ProductService{
		r: r,
	}
}

type ProductService struct {
	r repository.IProductRepository
}

func (r ProductService) ProductList() ([]Product.Product, error) {
	return r.r.ProductList(), nil
}