package config

//诸葛上报event_name
type zhuGeEventName struct {
	Qnr                    string
	UserLoginErr           string
	UserLoginSuc           string
	UserCertificateSendSuc string
	UserIdentify           string
	UserInvitedBy          string
}

var ZhuGeEventName = zhuGeEventName{
	Qnr:                    "绿色金融调查问卷", // 金融调查问卷
	UserLoginErr:           "用户登陆失败",   //用户登陆失败
	UserLoginSuc:           "用户登陆成功",   //用户登陆成功
	UserCertificateSendSuc: "携手同行-证书发放",
	UserIdentify:           "用户渠道",
	UserInvitedBy:          "用户邀请关系上报",
}
