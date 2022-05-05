package entity

import "mio/internal/pkg/model"

const (
	OrderByStepHistoryCountDesc OrderBy = "order_by_step_history_count_desc"
	OrderByStepHistoryTimeDesc  OrderBy = "order_by_step_history_time_desc"
)

type StepHistory struct {
	ID            int64      `json:"id"`
	UserId        int64      `json:"userId"`
	Count         int        `json:"count"`
	RecordedTime  model.Time `json:"recordedTime"`
	RecordedEpoch int64      `json:"recordedEpoch"`
}

func (StepHistory) TableName() string {
	return "step_history"
}
