package entity

import "time"

type CarbonScene struct {
	ID              int64                 `json:"id"`
	Pid             int64                 `json:"pid"`
	Type            CarbonTransactionType `json:"type"`
	Title           string                `json:"title"`
	Icon            string                `json:"icon"`
	Desc            string                `json:"desc"`
	UnitNumerator   float64               `json:"unitNumerator"`
	UnitDenominator float64               `json:"unitDenominator"`
	UnitDesc        string                `json:"unitDesc"`
	MaxCount        int                   `json:"maxCount"`
	MaxPoint        int                   `json:"maxPoint"`
	MaxCarbon       float64               `json:"maxCarbon"`
	Status          int8                  `json:"status"`
	CreatedAt       time.Time             `json:"createdAt"`
	UpdatedAt       time.Time             `json:"updatedAt"`
}
