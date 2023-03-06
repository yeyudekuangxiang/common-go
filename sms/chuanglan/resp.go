package chuanglan

type errorResp struct {
	Code   int    `json:"code"`
	Errmsg string `json:"errmsg"`
}

func (e errorResp) IsSuccess() bool {
	return e.Code == 0
}
