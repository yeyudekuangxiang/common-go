package activity

import "time"

type YtxLog struct {
	Id             int64     `json:"id"`
	OrderNo        string    `json:"orderNo"`
	OpenId         string    `json:"openId"`
	PlatformUserId string    `json:"platformUserId"`
	Amount         float64   `json:"amount"`
	AdditionalInfo string    `json:"additionalInfo"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func (YtxLog) TableName() string {
	return "activity_zyh_log"
}
