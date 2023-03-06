package chuanglan

type ErrorResp struct {
	Code   int    `json:"code"`
	Errmsg string `json:"errmsg"`
}

func (e ErrorResp) IsSuccess() bool {
	return e.Code == 0
}
