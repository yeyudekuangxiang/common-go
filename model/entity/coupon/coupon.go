package Coupon

type Coupon struct {
	ID           int    `gorm:"primary_key;column:id" json:"id"`
	CouponId     string `gorm:"column:coupon_id" json:"couponId"`
	CouponTypeId string `gorm:"column:coupon_type_id" json:"couponTypeId"`
	UpdateTime   string `gorm:"column:update_time" json:"updateTime"`
	Openid       string `gorm:"column:openid" json:"openid"`
	Redeemed     string `gorm:"column:redeemed" json:"redeemed"`
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
