package config

//诸葛上报event_name
type zhuGeEventName struct {
	Qnr                    string
	UserLoginErr           string
	UserLoginSuc           string
	UserCertificateSendSuc string
	UserIdentify           string
	UserInvitedBy          string
	MessageMiniSubscribe   string
	FirstIncPoint          string
	FirstDecPoint          string
	NewUserAdd             string
	MsgSecCheck            string
	MediaCheck             string
	DuiBaOrder             string
	YTXOrder               string
	YTXCollectPoint        string
	YTXReward              string
}

var ZhuGeEventName = zhuGeEventName{
	Qnr:                    "绿色金融调查问卷", // 金融调查问卷
	UserLoginErr:           "用户登陆失败",   //用户登陆失败
	UserLoginSuc:           "用户登陆成功",   //用户登陆成功
	UserCertificateSendSuc: "携手同行-证书发放",
	UserIdentify:           "用户渠道",
	UserInvitedBy:          "用户邀请关系上报",
	MessageMiniSubscribe:   "小程序订阅消息",
	FirstIncPoint:          "首次赚积分",
	FirstDecPoint:          "首次消耗积分",
	NewUserAdd:             "新用户",
	MsgSecCheck:            "文本内容审核",
	MediaCheck:             "媒体文件审核",
	DuiBaOrder:             "兑吧商城订单",
	YTXReward:              "亿通行-发放奖励",
	YTXOrder:               "亿通行-完成乘车",
	YTXCollectPoint:        "亿通行-收取气泡",
}
