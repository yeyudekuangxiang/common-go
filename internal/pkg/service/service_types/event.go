package service_types

type SubmitOrderForEventParam struct {
	UserId  int64
	EventId string
}

type SubmitOrderForEventResult struct {
	CertificateNo string `json:"certificateNo"`
}
