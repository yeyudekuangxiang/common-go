package message

type SendData struct {
	BizId            string `json:"bizId"`
	Body             string `json:"body"`
	ContentType      string `json:"contentType"`
	Method           string `json:"method"`
	SuccessHttpCodes []int  `json:"successHttpCodes"`
	Token            string `json:"token"`
	Url              string `json:"url"`
}
