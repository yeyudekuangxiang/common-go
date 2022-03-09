package Product

type Product struct {
	ID                     int    `gorm:"primary_key;column:id" json:"id"`
	ProductItemId          string `gorm:"column:product_item_id" json:"productItemId"`
	Virtual                string `gorm:"column:virtual" json:"virtual"`
	Title                  string `gorm:"column:title" json:"title"`
	Cost                   string `gorm:"column:cost" json:"cost"`
	ImageUrl               string `gorm:"column:image_url" json:"imageUrl"`
	RemainingCount         string `gorm:"column:remaining_count" json:"remainingCount"`
	SalesCount             string `gorm:"column:sales_count" json:"salesCount"`
	Active                 string `gorm:"column:active" json:"active"`
	ProductItemReferenceId string `gorm:"column:product_item_reference_id" json:"productItemReferenceId"`
	Sort                   string `gorm:"column:sort" json:"sort"`
}
