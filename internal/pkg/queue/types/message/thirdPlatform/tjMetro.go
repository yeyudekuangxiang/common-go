package thirdPlatform

type TjMetroMessage struct {
	OpenId           string
	ThirdCouponTypes int64
}

type TjMetroSendMessage struct {
	Uid                 int64
	OpenId              string
	Phone               string
	BizId               string
	CouponCardTypeId    int64
	DistributionChannel string
}
