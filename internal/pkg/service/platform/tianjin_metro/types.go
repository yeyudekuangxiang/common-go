package tianjin_metro

type SyncMetroRequest struct {
	OpenId         string `json:"openId"`              //亿通行openId
	RegDate        string `json:"regDate"`             //注册时间，格式yyyyMMddHHmmss
	PlatformUserId string `json:"platformUserId"`      //绿喵用户ID
	Ts             int64  `json:"ts"`                  //时间戳，毫秒
	Signature      string `json:"signature,omitempty"` //签名，计算获得
}

type MetroRequest struct {
	AllotId     string `json:"allotId"`
	EtUserPhone string `json:"etUserPhone"`
	AllotNum    int8   `json:"allotNum"`
}

type MetroResponse struct {
	SubCode    string       `json:"subCode,omitempty"`
	SubData    MetroSubData `json:"subData"`
	SubMessage string       `json:"subMessage,omitempty"`
}

type MetroSubData struct {
	OrderNo     string `json:"orderNo"`
	TradeNo     string `json:"tradeNo,omitempty"`
	Success     bool   `json:"success"`
	Description string `json:"description,omitempty"`
}
