package entity

import (
	"mio/internal/pkg/model"
)

type PointTransactionCountLimit struct {
	Id              int64                `json:"id"`
	OpenId          string               `json:"openId"`
	TransactionType PointTransactionType `json:"transactionType"`
	MaxCount        int                  `json:"maxCount"`
	CurrentCount    int                  `json:"currentCount"`
	UpdateTime      model.Time           `json:"updateTime"`
	TransactionDate model.Date           `json:"transactionDate"`
}
