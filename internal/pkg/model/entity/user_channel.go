package entity

import "mio/internal/pkg/model"

type UserChannel struct {
	ID         int64      `gorm:"primary_key;column:id" json:"id"`
	Cid        int64      `gorm:"column:cid" json:"cid"`
	Pid        int64      `gorm:"column:pid" json:"pid"`
	Name       string     `gorm:"column:name" json:"name"`
	Code       string     `gorm:"column:code" json:"code"`
	Company    string     `gorm:"column:company" json:"company"`
	CreateTime model.Time `json:"createTime"`
	UpdateTime model.Time `gorm:"column:update_time" json:"updateTime"`
}

func (UserChannel) TableName() string {
	return "user_channel"
}
