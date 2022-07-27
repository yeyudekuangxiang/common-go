package coupon

import (
	"mio/internal/pkg/util/timeutils"
)

type CouponRecord struct {
	ID               int64
	Openid           string
	UpdateTime       timeutils.Time
	Name             string
	CouponId         string
	CouponTypeId     string
	ProductItemId    string
	ProductItemCount int
	Point            int
}
