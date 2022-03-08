package Product

type Product struct {
	ID                     int    `gorm:"primary_key;column:id" json:"id"`
	ProductItemId          string `gorm:"column:product_item_id" json:"productItemId"`
	Virtual                string `gorm:"column:virtual" json:"virtual"`
	Title                  string `gorm:"title" json:"title"`
	Cost                   string `gorm:"cost" json:"cost"`
	ImageUrl               string `gorm:"image_url" json:"imageUrl"`
	RemainingCount         string `gorm:"remaining_count" json:"remainingCount"`
	SalesCount             string `gorm:"sales_count" json:"salesCount"`
	Active                 string `gorm:"active" json:"active"`
	ProductItemReferenceId string `gorm:"product_item_reference_id" json:"productItemReferenceId"`
	Sort                   string `gorm:"sort" json:"sort"`
}
