package jhx

type CommonErr struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (c *CommonErr) IsSuccess() bool {
	return c.Code == 0
}

type TicketCreateResponse struct {
	CommonErr
	Time string           `json:"time"`
	Data TicketCreateData `json:"data"`
}

type TicketCreateData struct {
	QrCodeStr  string `json:"qrcodestr" form:"qrcodestr"`
	ExpireTime string `json:"expireTime" from:"expireTime"`
}
