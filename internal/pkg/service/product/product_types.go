package product

type GetProductItemListParam struct {
	ItemIds []string
}
type CreateOrUpdateProductItemParam struct {
	ItemId   string
	Virtual  bool
	Title    string
	Cost     int
	ImageUrl string
	Sort     int
}
