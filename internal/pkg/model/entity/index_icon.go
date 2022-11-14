package entity

import (
	"time"
)

type IndexIcon struct {
	ID        int64           `gorm:"primary_key;column:id" json:"id"`
	Title     string          `gorm:"column:title" json:"title"`
	Type      string          `gorm:"column:type" json:"type"`
	RowNum    string          `gorm:"column:row_num" json:"rowNum"`
	Sort      int8            `gorm:"column:sort" json:"sort"`
	Status    IndexIconStatus `gorm:"column:status" json:"status"`
	IsOpen    int8            `gorm:"column:is_open" json:"isOpen"`
	Pic       string          `gorm:"column:pic" json:"pic"`
	CreatedAt time.Time       `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time       `gorm:"column:updated_at" json:"updatedAt"`
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

func (p IndexIconStatus) Text() string {
	switch p {
	case IndexIconStatusOk:
		return "上线"
	case IndexIconStatusDown:
		return "下线"
	}
	return "未知"
}

type IndexIconIsOpen int8

const (
	IndexIconIsOpenOk   IndexIconIsOpen = 1 //上线
	IndexIconIsOpenDown IndexIconIsOpen = 2 //下线
)

var (
	IndexIconIsOpenMap = map[IndexIconIsOpen]string{
		IndexIconIsOpenOk:   "开启",
		IndexIconIsOpenDown: "关闭",
	}
)

func (p IndexIconIsOpen) Text() string {
	switch p {
	case IndexIconIsOpenOk:
		return "开启"
	case IndexIconIsOpenDown:
		return "关闭"
	}
	return "未知"
}
