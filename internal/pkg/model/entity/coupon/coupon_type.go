package coupon

import "mio/internal/pkg/model"

type CouponType struct {
	ID               int64      `json:"id"`
	Name             string     `json:"name"`
	CouponTypeId     string     `json:"couponTypeId"`
	CreateTime       model.Time `json:"createTime"`
	StartTime        model.Time `json:"startTime"`
	EndTime          model.Time `json:"endTime"`
	Point            int        `json:"point"`
	ProductItemId    string     `json:"productItemId"`
	ProductItemCount int        `json:"productItemCount"`
	Partnership      string     `json:"partnership"`
	Description      string     `json:"description"`
	External         bool       `json:"external"`
	Active           bool       `json:"active"`
	Reusable         bool       `json:"reusable"`
}
