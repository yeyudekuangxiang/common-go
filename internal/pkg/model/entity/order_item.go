package entity

type OrderItem struct {
	Id      int64  `gorm:"primary_key" json:"id"`
	OrderId string `json:"orderId"`
	ItemId  string `json:"itemId"`
	Count   int    `json:"itemCount"`
	Cost    int    `json:"cost"`
}
