package jhx

type CommonResponse struct {
	Code int                  `json:"code"`
	Msg  string               `json:"msg"`
	Time string               `json:"time"`
	Data TicketCreateResponse `json:"data"`
}

type TicketCreateResponse struct {
	QrCodeStr  string `json:"qrcodestr" form:"qrcodestr"`
	ExpireTime string `json:"expireTime" from:"expireTime"`
}
