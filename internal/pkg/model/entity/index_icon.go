package entity

import (
	"mio/internal/pkg/model"
)

type IndexIcon struct {
	ID        int64      `gorm:"primary_key;column:id" json:"id"`
	Title     string     `gorm:"column:title" json:"title"`
	Type      string     `gorm:"column:type" json:"type"`
	RowNum    string     `gorm:"column:row_num" json:"rowNum"`
	Sort      int8       `gorm:"column:sort" json:"sort"`
	Status    int8       `gorm:"column:status" json:"status"`
	IsOpen    int8       `gorm:"column:is_open" json:"isOpen"`
	Pic       string     `gorm:"column:pic" json:"pic"`
	CreatedAt model.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt model.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (IndexIcon) TableName() string {
	return "index_icon"
}

type IndexIconStatus int8

const (
	IndexIconStatusOk   IndexIconStatus = 1 //上线
	IndexIconStatusDown IndexIconStatus = 2 //下线
)

var (
	IndexIconStatusMap = map[IndexIconStatus]string{
		IndexIconStatusOk:   "上线",
		IndexIconStatusDown: "下线",
	}
)
