package config

//诸葛上报event_name
type zhuGeEventName struct {
	Qnr                    string
	UserLoginErr           string
	UserLoginSuc           string
	UserCertificateSendErr string
	UserCertificateSendSuc string
	UserIdentify           string
}

var ZhuGeEventName = zhuGeEventName{
	Qnr:                    "绿色金融调查问卷", // 金融调查问卷
	UserLoginErr:           "用户登陆失败",   //用户登陆失败
	UserLoginSuc:           "用户登陆成功",   //用户登陆成功
	UserCertificateSendErr: "携手同行-证书-发放成功",
	UserCertificateSendSuc: "携手同行-证书-发放失败",
	UserIdentify:           "用户属性",
}
