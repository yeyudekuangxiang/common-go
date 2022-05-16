package entity

import "mio/internal/pkg/model"

type Badge struct {
	ID            int64      `json:"id"`
	Code          string     `json:"code"`
	OpenId        string     `json:"openId" gorm:"column:openid"`
	CertificateId string     `json:"certificateId"`
	ProductItemId string     `json:"productItemId"`
	CreateTime    model.Time `json:"createTime"`
	Partnership   string     `json:"partnership"`
	OrderId       string     `json:"orderId"`
}
