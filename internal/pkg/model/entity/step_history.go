package entity

import "mio/internal/pkg/model"

const (
	OrderByStepHistoryCountDesc OrderBy = "order_by_step_history_count_desc"
)

type StepHistory struct {
	ID            int64      `json:"id"`
	OpenId        string     `json:"openid"`
	Count         int        `json:"count"`
	RecordedTime  model.Date `json:"recordedTime"`
	RecordedEpoch int64      `json:"recordedEpoch"`
}

func (StepHistory) TableName() string {
	return "step_history"
}
