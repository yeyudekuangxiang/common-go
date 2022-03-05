package product

import (
	"github.com/gin-gonic/gin"
	"mio/service"
)

var DefaultProductController = ProductController{}

type ProductController struct {
}

func (ProductController) ProductList(c *gin.Context) (gin.H, error) {
	list, err := service.DefaultProductService.ProductList()

	return gin.H{
		"records":          list,
		"total":            76,
		"size":             100,
		"current":          1,
		"orders":           nil,
		"optimizeCountSql": true,
		"hitCount":         false,
		"countId":          nil,
		"maxLimit":         nil,
		"searchCount":      true,
		"pages":            1,
	}, err
}
