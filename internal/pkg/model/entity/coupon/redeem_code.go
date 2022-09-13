package coupon

import "time"

type RedeemCode struct {
	ID         int64
	CodeId     string
	CouponId   string
	CreateTime time.Time
	UpdateTime time.Time
}
