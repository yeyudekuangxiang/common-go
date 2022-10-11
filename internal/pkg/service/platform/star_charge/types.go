package star_charge

type starChargeResponse struct {
	Ret  int    `json:"Ret"`
	Msg  string `json:"Msg"`
	Data string `json:"Data"`
	Sig  string `json:"Sig"`
}

type starChargeAccessResult struct {
	OperatorID         string `json:"operatorID,omitempty"`
	SucStat            int    `json:"sucStat,omitempty"`
	AccessToken        string `json:"accessToken,omitempty"`
	TokenAvailableTime int    `json:"tokenAvailableTime,omitempty"`
	FailReason         int    `json:"failReason,omitempty"`
}

type starChargeProvideResult struct {
	SuccStat      int    `json:"succStat,omitempty"`
	FailReason    int    `json:"failReason,omitempty"`
	FailReasonMsg string `json:"failReasonMsg,omitempty"`
	CouponCode    string `json:"couponCode,omitempty"`
}
