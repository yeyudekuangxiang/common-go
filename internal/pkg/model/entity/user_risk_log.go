package entity

import "mio/internal/pkg/model"

type UserRiskLog struct {
	ID        int
	OpenId    string `json:"openid"`
	Scene     int64  `json:"scene"`
	MobileNo  string `json:"mobile_no"`
	ClientIp  string `json:"client_ip"`
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	UnoinId   int    `json:"unoin_id"`
	RiskRank  int    `json:"risk_rank"`
	CreatedAt model.Time
}

func (UserRiskLog) TableName() string {
	return "user_risk_log"
}
