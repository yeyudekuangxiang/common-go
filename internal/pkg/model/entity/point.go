package entity

type Point struct {
	Id      int    `json:"column:id"`
	OpenId  string `gorm:"column:openid"`
	Balance int64  `json:"column:balance"`
}

func (Point) TableName() string {
	return "point"
}
