package entity

import "time"

type UserExtendInfo struct {
	ID           int
	Openid       string    `json:"openid"`
	Uid          int64     `json:"uid"`
	Ip           string    `json:"ip"`
	Province     string    `json:"province"`
	City         string    `json:"city"`
	District     string    `json:"district"`
	StreetNumber string    `json:"street_number"`
	Adcode       string    `json:"adcode"`
	Street       string    `json:"street"`
	CityCode     int       `json:"city_code"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (UserExtendInfo) TableName() string {
	return "user_extend_info"
}
