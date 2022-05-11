package wxapp

type Response struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type UnlimitedQRCodeResponse struct {
	Response
	ContentType string `json:"contentType"`
	Buffer      []byte `json:"buffer"`
}
type UserRiskRankParam struct {
	AppId        string `json:"appid"`
	OpenId       string `json:"openid"`
	Scene        int64  `json:"scene"`
	MobileNo     string `json:"mobile_no"`
	ClientIp     string `json:"client_ip"`
	EmailAddress string `json:"email_address"`
	ExtendedInfo string `json:"extended_info"`
	IsTest       bool   `json:"is_test"`
}
type UserRiskRankResponse struct {
	Response
	UnoinId  int `json:"unoin_id"`
	RiskRank int `json:"risk_rank"`
}
