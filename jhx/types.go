package jhx

type commonResponse struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Time string                 `json:"time"`
	Data map[string]interface{} `json:"data"`
}

type TicketCreateResponse struct {
	QrCodeStr  string `json:"qrcodestr" form:"qrcodestr"`
	ExpireTime string `json:"expireTime" from:"expireTime"`
}
