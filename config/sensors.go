package config

var SensorsEventName = sensorsEventName{
	Qnr:                  "GreenFinanceQnr", // 金融调查问卷
	UserInvitedBy:        "UserInvite",
	MessageMiniSubscribe: "小程序订阅消息",
	FirstIncPoint:        "首次赚积分",
	FirstDecPoint:        "首次消耗积分",
	NewUserAdd:           "NewUser",
	MsgSecCheck:          "MsgSecCheck",
	DuiBaOrder:           "GreenMall",
	YTX:                  "YtxActivity",
	CommunityTopic:       "CommunityTopic",
	PointChange:          "PointChange",
}

//诸葛上报event_name
type sensorsEventName struct {
	Qnr                  string
	UserInvitedBy        string
	MessageMiniSubscribe string
	FirstIncPoint        string
	FirstDecPoint        string
	NewUserAdd           string
	MsgSecCheck          string
	DuiBaOrder           string
	YTX                  string
	CommunityTopic       string
	PointChange          string
}
