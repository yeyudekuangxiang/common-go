package config

//诸葛上报event_name
type message struct {
	Qnr                    string
	UserLoginErr           string
	UserLoginSuc           string
	UserCertificateSendSuc string
	UserIdentify           string
}

var messageStruct = message{
	Qnr:                    "绿色金融调查问卷", // 金融调查问卷
	UserLoginErr:           "用户登陆失败",   //用户登陆失败
	UserLoginSuc:           "用户登陆成功",   //用户登陆成功
	UserCertificateSendSuc: "携手同行-证书发放",
	UserIdentify:           "用户渠道",
}

type MessageDataValueStruct struct {
	Value string `json:"value"`
}

type MessageDataStruct struct {
	Number1 MessageDataValueStruct `json:"number01"`
	Thing5  MessageDataValueStruct `json:"date01"`
	Time2   MessageDataValueStruct `json:"site01"`
	Number3 MessageDataValueStruct `json:"site02"`
}

type MessageStruct struct {
	Touser           string            `json:"touser"`
	TemplateId       string            `json:"template_id"`
	Page             string            `json:"page"`
	MiniprogramState string            `json:"miniprogram_state"`
	Lang             string            `json:"lang"`
	Data             MessageDataStruct `json:"data"`
}
