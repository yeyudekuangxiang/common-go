package srv_types

type SubmitOrderForEventParam struct {
	UserId  int64
	EventId string
}

type SubmitOrderForEventResult struct {
	CertificateNo string `json:"certificateNo"`
	UploadCode    string `json:"uploadCode"`
}

type SubmitOrderForEventGDParam struct {
	UserId              int64
	EventId             string
	WechatServiceOpenId string
	OpenId              string
}
