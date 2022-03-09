package Pugc

import "time"

type Pugc struct {
	ID          int       `gorm:"primary_key;column:id" json:"id"`
	UserId      int       `gorm:"column:user_id" json:"user_id"`
	Title       string    `gorm:"column:title" json:"title"`
	Content     string    `gorm:"column:content" json:"content"`
	Pic         string    `gorm:"column:pic" json:"pic"`
	CreatedTime time.Time `gorm:"column:created_time" json:"created_time"`
}

type PugcAddModel struct {
	UserId      int       `gorm:"column:user_id" json:"user_id"`
	Title       string    `gorm:"column:title" json:"title"`
	Content     string    `gorm:"column:content" json:"content"`
	Pic         string    `gorm:"column:pic" json:"pic"`
	CreatedTime time.Time `gorm:"column:created_time" json:"created_time"`
}
