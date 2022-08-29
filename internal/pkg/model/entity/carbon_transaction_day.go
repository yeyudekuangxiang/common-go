package entity

import (
	"time"
)

type CarbonTransactionDay struct {
	ID        int64     `json:"id"`
	OpenId    string    `gorm:"column:openid" json:"openId"`
	UserId    int64     `json:"userId"`
	VDate     time.Time `json:"vDate"`
	Value     float64   `json:"value"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
