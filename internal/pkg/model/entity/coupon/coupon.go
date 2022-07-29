package coupon

import "mio/internal/pkg/model"

type Coupon struct {
	ID               int        `gorm:"primary_key;column:id" json:"id"`
	CouponId         string     `gorm:"column:coupon_id" json:"couponId"`
	CouponTypeId     string     `gorm:"column:coupon_type_id" json:"couponTypeId"`
	CreateTime       model.Time `json:"createTime"`
	UpdateTime       model.Time `gorm:"column:update_time" json:"updateTime"`
	StartTime        model.Time `json:"startTime"`
	EndTime          model.Time `json:"endTime"`
	Openid           string     `gorm:"column:openid" json:"openid"`
	Redeemed         bool       `gorm:"column:redeemed" json:"redeemed"`
	OrderReferenceId string     `json:"orderReference"`
}

type CouponRes struct {
	ID           int    `gorm:"primary_key;column:id" json:"id"`
	CouponId     string `gorm:"column:coupon_id" json:"couponId"`
	CouponTypeId string `gorm:"column:coupon_type_id" json:"couponTypeId"`
	CouponType   string `gorm:"column:coupon_type" json:"couponType"`
	UpdateTime   string `gorm:"column:update_time" json:"updateTime"`
	Openid       string `gorm:"column:openid" json:"openid"`
	Redeemed     string `gorm:"column:redeemed" json:"redeemed"`
}
