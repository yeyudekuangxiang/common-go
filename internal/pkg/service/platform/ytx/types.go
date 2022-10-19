package ytx

type SynchroRequest struct {
	OpenId         string `json:"openId"`             //亿通行openId
	RegDate        string `json:"regDate"`            //注册时间，格式yyyyMMddHHmmss
	PlatformUserId string `json:"platformUserId"`     //绿喵用户ID
	Ts             int64  `json:"ts"`                 //时间戳，毫秒
	Signature      string `json:"platform,omitempty"` //签名，计算获得
}

type synchroResponse struct {
	ResCode    string                 `json:"resCode"`    //返回码
	ResMessage string                 `json:"resMessage"` //返回描述
	ResData    map[string]interface{} `json:"resData"`    //object
}

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

type GrantV2Response struct {
	SubCode    string         `json:"subCode,omitempty"`
	SubData    GrantV2SubData `json:"subData"`
	SubMessage string         `json:"subMessage,omitempty"`
}

type GrantV2SubData struct {
	OrderNo     string `json:"orderNo"`
	TradeNo     string `json:"tradeNo,omitempty"`
	Success     bool   `json:"success"`
	Description string `json:"description,omitempty"`
}
