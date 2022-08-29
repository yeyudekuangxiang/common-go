package config

//诸葛上报event_name
type zhuGeEventName struct {
	Qnr string
}

var ZhuGeEventName = zhuGeEventName{
	Qnr: "zhu_ge_qnr_green_finance", // 金融调查问卷
}
