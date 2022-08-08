package entity

import (
	"time"
)

type CarbonTransactionCountLimitDay struct {
	ID           int64                 `json:"id"`
	OpenId       string                `gorm:"column:openid" json:"openId"`
	UserId       int64                 `json:"userId"`
	Type         CarbonTransactionType `json:"type"`
	MaxCount     int                   `json:"maxCount"`
	CurrentCount int                   `json:"currentCount"`
	VDate        time.Time             `json:"vDate"`
	CreatedAt    time.Time             `json:"createdAt"`
	UpdatedAt    time.Time             `json:"updatedAt"`
}
