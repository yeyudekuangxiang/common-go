package entity

type Carbon struct {
	Id     int     `json:"column:id"`
	OpenId string  `gorm:"column:openid"`
	Carbon float64 `json:"column:carbon"`
}

func (Carbon) TableName() string {
	return "carbon"
}
