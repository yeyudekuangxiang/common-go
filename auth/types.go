package auth

type UserIdentityVerificationReq struct {
	Name          string `json:"name"`
	IdentityCard  string `json:"identityCard"`
	Phone         string `json:"phone"`
	transactionId string `json:"transactionId"`
}

type UserIdentityVerificationResp struct {
	Status int `json:"status"`
	Result struct {
		AuthResult bool   `json:"authResult"`
		Msg        string `json:"msg"`
		AuthnBizNo string `json:"authnBizNo"`
	} `json:"result"`
	Error     string `json:"error"`
	ErrorCode string `json:"errorCode"`
	Msg       string `json:"msg"`
	Time      string `json:"time"`
}
