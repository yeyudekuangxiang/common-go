package config

var SensorsEventName = sensorsEventName{
	Qnr:                  "GreenFinanceQnr", // 金融调查问卷
	MessageMiniSubscribe: "小程序订阅消息",
	MsgSecCheck:          "MsgSecCheck",
	DuiBaOrder:           "GreenMall",
	YTX:                  "YtxActivity",
	CommunityTopic:       "CommunityTopic",
	ActivityApply:        "ActivityApply",
}

//诸葛上报event_name
type sensorsEventName struct {
	Qnr                  string
	MessageMiniSubscribe string
	MsgSecCheck          string
	DuiBaOrder           string
	YTX                  string
	CommunityTopic       string
	ActivityApply        string
}
