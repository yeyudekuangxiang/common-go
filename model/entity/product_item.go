package entity

type ProductItem struct {
	ID                     int64
	ProductItemId          string
	Virtual                bool
	Title                  string
	Cost                   int
	ImageUrl               string
	RemainingCount         int  //库存
	SalesCount             int  //已售
	Active                 bool //是否上架
	ProductItemReferenceId string
	Sort                   int
}
