package entity

import (
	"time"
)

type City struct {
	ID        int64
	CityCode  string
	Name      string
	PidCode   string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
