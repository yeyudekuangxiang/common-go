package entity

import "time"

type ExampleTest struct {
	ID        int64     `gorm:"primary_key;column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`

	Title  string `gorm:"column:title" json:"title"`
	Type   string `gorm:"column:type" json:"type"`
	RowNum string `gorm:"column:row_num" json:"rowNum"`
	Sort   int8   `gorm:"column:sort" json:"sort"`
	Status int8   `gorm:"column:status" json:"status"`
	IsOpen int8   `gorm:"column:is_open" json:"isOpen"`
	Pic    string `gorm:"column:pic" json:"pic"`
}
