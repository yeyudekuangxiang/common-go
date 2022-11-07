package entity

import (
	"time"
)

type CouponHistory struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OpenId       string    `gorm:"type:varchar(200)" json:"openId"`
	CouponType   string    `gorm:"type:varchar(50)" json:"couponType"`
	Code         string    `gorm:"type:varchar(100)" json:"code"`
	AdditionInfo string    `json:"additionInfo"`
	CreateTime   time.Time `json:"createTime"`
	UpdateTime   time.Time `json:"updateTime"`
}

func (CouponHistory) TableName() string {
	return "coupon_history"
}
