package ytx

type GrantV2Request struct {
	AppId       string         `json:"appId"`
	AppSecret   string         `json:"appSecret"`
	Ts          string         `json:"ts"`
	Version     string         `json:"version,omitempty"`
	AppVersion  string         `json:"appVersion,omitempty"`
	DeviceType  string         `json:"deviceType,omitempty"`
	DeviceId    string         `json:"deviceId,omitempty"`
	RedirectUrl string         `json:"redirectUrl,omitempty"`
	OpenId      string         `json:"openId,omitempty"`
	UserId      string         `json:"userId,omitempty"`
	ReqData     GrantV2ReqData `json:"reqData"`
}

type GrantV2ReqData struct {
	OrderNo  string  `json:"orderNo"`
	PoolCode string  `json:"poolCode"`
	Amount   float64 `json:"amount"`
	UserId   string  `json:"userId,omitempty"`
	PhoneNum string  `json:"phoneNum,omitempty"`
	OpenId   string  `json:"openId,omitempty"`
	Remark   string  `json:"remark"`
}

type CommonErr struct {
	SubCode    string `json:"subCode,omitempty"`
	SubMessage string `json:"subMessage,omitempty"`
}

func (c *CommonErr) IsSuccess() bool {
	return c.SubCode == "0000"
}

type GrantV2Response struct {
	CommonErr
	SubData GrantV2SubData `json:"subData"`
}

type GrantV2SubData struct {
	OrderNo     string `json:"orderNo"`
	TradeNo     string `json:"tradeNo,omitempty"`
	Success     bool   `json:"success"`
	Description string `json:"description,omitempty"`
}
