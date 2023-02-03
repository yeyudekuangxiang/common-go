package smsmsg

import (
	"encoding/json"
	"strings"
)

type MsgMessage struct {
	Phone string `json:"phone"`
	Msg   string `json:"msg"`
}

type SmsMessage struct {
	TemplateKey string `json:"templateKey"`
	Phone       string `json:"phone"`
	Args        Args   `json:"msg"`
}

func (h SmsMessage) Byte() ([]byte, error) {
	return json.Marshal(h)
}

type HttpSmsMessage struct {
	Url              string            `json:"url"`
	Token            string            `json:"token"`
	Method           string            `json:"method"`
	ContentType      string            `json:"contentType"`
	Body             string            `json:"body"`
	Form             map[string]string `json:"form"`
	Header           map[string]string `json:"header"`
	SuccessHttpCodes []int             `json:"successHttpCodes"`
}

func (h HttpSmsMessage) Byte() ([]byte, error) {
	return json.Marshal(h)
}

type Args []string

func (a Args) Byte() ([]byte, error) {
	return []byte("\"" + strings.Join(a, ",") + "\""), nil
}
