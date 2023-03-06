package chuanglan

type ErrorResp struct {
	Code   string `json:"code"`
	Errmsg string `json:"errorMsg"`
}

func (e ErrorResp) IsSuccess() bool {
	return e.Code == "0"
}
