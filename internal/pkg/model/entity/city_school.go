package entity

import "time"

type CitySchool struct {
	Id         int64     `json:"id"`
	CityCode   string    `json:"cityCode"`
	SchoolName string    `json:"schoolName"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (s CitySchool) TableName() string {
	return "city_school"
}
