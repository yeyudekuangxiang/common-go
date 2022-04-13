package entity

import "mio/internal/pkg/model"

type Session struct {
	ID               int64      `gorm:"id"`
	OpenId           string     `gorm:"openid"`
	SessionKey       string     `gorm:"session_key"`
	WechatSessionKey string     `gorm:"wechat_session_key"`
	Time             model.Time `gorm:"time"`
	WxUnionId        string     `gorm:"wx_union_id"`
}
