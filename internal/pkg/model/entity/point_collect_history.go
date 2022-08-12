package entity

import "mio/internal/pkg/model"

type PointCollectHistory struct {
	ID     int64      `json:"id"`
	OpenId string     `json:"openId" gorm:"column:openid"`
	Type   string     `json:"type"`
	Info   string     `json:"info"`
	Date   model.Date `json:"date"`
	Time   model.Time `json:"time"`
}

type PointCollectLog struct {
	ID              int64      `json:"id"`
	OpenId          string     `json:"openId" gorm:"column:openid"`
	Type            string     `json:"type"`
	Point           int64      `json:"point"`
	AdditionalOrder string     `json:"additionalOrder"`
	Info            string     `json:"info"`
	Date            model.Date `json:"date"`
	Time            model.Time `json:"time"`
}
