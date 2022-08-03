package entity

import "mio/internal/pkg/model"

type UserChannel struct {
	ID         int64      `gorm:"primary_key;column:id" json:"id"`
	Cid        int64      `gorm:"column:cid" json:"cid"`         //渠道id
	Pid        int64      `gorm:"column:pid" json:"pid"`         //父id
	Name       string     `gorm:"column:name" json:"name"`       //渠道名
	Code       string     `gorm:"column:code" json:"code"`       //渠道编号
	Company    string     `gorm:"column:company" json:"company"` //公司名
	CreateTime model.Time `json:"createTime"`
	UpdateTime model.Time `gorm:"column:update_time" json:"updateTime"`
}

func (UserChannel) TableName() string {
	return "user_channel"
}
