package entity

import "mio/internal/pkg/model"

const (
	OrderByCheckinHistoryTimeDesc OrderBy = "order_by_checkin_history_time_desc"
)

type CheckinHistory struct {
	Id            int64
	OpenId        string
	CheckedNumber int
	Time          model.Time
}

func (CheckinHistory) TableName() string {
	return "check_in_history"
}
