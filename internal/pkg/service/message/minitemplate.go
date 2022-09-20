package message

import (
	"github.com/medivhzhan/weapp/v3/subscribemessage"
	"mio/config"
)

type IMiniSubTemplate interface {
	ToData() map[string]subscribemessage.SendValue
	TemplateId() string
	IDIUDGIUGISIAHSUIAHUISAHUIAGUISGIU()
}

//MiniPointSendTemplate 积分变更
type MiniPointSendTemplate struct {
	Point    string `json:"point"`
	Source   string `json:"source"`
	Time     string `json:"time"`
	AllPoint string `json:"tip"`
}

func (m MiniPointSendTemplate) ToData() map[string]subscribemessage.SendValue {
	return map[string]subscribemessage.SendValue{
		"number1": {Value: m.Point}, //strconv.Itoa(m.Money)strconv.Atoi(m.Point)
		"thing5":  {Value: m.Source},
		"time2":   {Value: m.Time},
		"number3": {Value: m.AllPoint},
	}
}

func (m MiniPointSendTemplate) TemplateId() string {
	//return "hQzUwkGMYqgNsKOad7RnIwwBpfkVfsuJvW6UqymwI8k"
	return config.MessageTemplateIds.SendPoint
}

func (m MiniPointSendTemplate) IDIUDGIUGISIAHSUIAHUISAHUIAGUISGIU() {
	return
}
