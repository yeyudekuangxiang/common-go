package coupon

import "gitlab.miotech.com/miotech-application/backend/common-go/tool/timetool"

type CouponRecord struct {
	ID               int64
	Openid           string
	UpdateTime       timetool.Time
	Name             string
	CouponId         string
	CouponTypeId     string
	ProductItemId    string
	ProductItemCount int
	Point            int
}
