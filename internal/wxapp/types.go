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
