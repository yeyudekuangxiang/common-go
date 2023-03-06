package weixinmobile

type errorResp struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func (e errorResp) IsSuccess() bool {
	return e.Errcode == 0
}
