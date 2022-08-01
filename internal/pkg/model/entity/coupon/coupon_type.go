package coupon

import "mio/internal/pkg/model"

type CouponType struct {
	ID               int64            `json:"id"`
	Name             string           `json:"name"`
	CouponTypeId     string           `json:"couponTypeId"`
	CreateTime       model.Time       `json:"createTime"`
	StartTime        model.Time       `json:"startTime"`
	EndTime          model.Time       `json:"endTime"`
	Point            model.NullInt    `json:"point"`
	ProductItemId    model.NullString `json:"productItemId"`
	ProductItemCount model.NullInt    `json:"productItemCount"`
	Partnership      model.NullString `json:"partnership"`
	Description      model.NullString `json:"description"`
	External         bool             `json:"external"`
	Active           bool             `json:"active"`
	Reusable         bool             `json:"reusable"`
}
