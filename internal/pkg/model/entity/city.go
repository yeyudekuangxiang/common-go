package entity

import (
	"time"
)

type City struct {
	ID        int64     `json:"id"`
	CityCode  string    `json:"cityCode"`
	Name      string    `json:"name"`
	PidCode   string    `json:"pidCode"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (c City) TableName() string {
	return "city"
}
