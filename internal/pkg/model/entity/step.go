package entity

import "mio/internal/pkg/model"

type Step struct {
	ID             int64
	OpenId         string
	Total          int64
	LastCheckTime  model.Time
	LastCheckCount int
}

func (Step) TableName() string {
	return "step"
}
