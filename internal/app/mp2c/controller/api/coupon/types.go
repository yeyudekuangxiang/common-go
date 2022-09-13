package coupon

type RedeemCodeForm struct {
	RedeemCodeId string `json:"redeemCodeId" form:"redeemCodeId" binding:"required"`
}
