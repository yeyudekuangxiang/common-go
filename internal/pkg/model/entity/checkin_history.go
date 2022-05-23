package entity

import "mio/internal/pkg/model"

const (
	OrderByCheckinHistoryTimeDesc OrderBy = "order_by_checkin_history_time_desc"
)

type CheckinHistory struct {
	Id            int64      `json:"id"`
	OpenId        string     `gorm:"column:openid"`
	CheckedNumber int        `json:"checkedNumber"`
	Time          model.Time `json:"time"`
}

func (CheckinHistory) TableName() string {
	return "check_in_history"
}
