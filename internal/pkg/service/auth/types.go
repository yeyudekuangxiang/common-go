package auth

type OaSignResult struct {
	AppId     string `json:"appId"`
	Timestamp string `json:"timestamp"`
	NonceStr  string `json:"nonceStr"`
	Signature string `json:"signature"`
}

type FindOaAuthWhiteBy struct {
	Domain string
	AppId  string
}
