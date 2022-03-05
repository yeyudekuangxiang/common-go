package Coupon

type Coupon struct {
	ID           int    `gorm:"primary_key;column:id" json:"id"`
	CouponId     string `gorm:"coupon_id" json:"couponId"`
	CouponTypeId string `gorm:"coupon_type_id" json:"couponTypeId"`
	UpdateTime   string `gorm:"update_time" json:"updateTime"`
	Openid       string `gorm:"openid" json:"openid"`
	Redeemed     string `gorm:"redeemed" json:"redeemed"`
}

type CouponRes struct {
	ID           int    `gorm:"primary_key;column:id" json:"id"`
	CouponId     string `gorm:"coupon_id" json:"couponId"`
	CouponTypeId string `gorm:"coupon_type_id" json:"couponTypeId"`
	CouponType   string `gorm:"coupon_type" json:"couponType"`
	UpdateTime   string `gorm:"update_time" json:"updateTime"`
	Openid       string `gorm:"openid" json:"openid"`
	Redeemed     string `gorm:"redeemed" json:"redeemed"`
}
