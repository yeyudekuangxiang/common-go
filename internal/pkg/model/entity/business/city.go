package business

import (
	"mio/internal/pkg/model"
	"time"
)

type AreaLevel string

const (
	AreaProvince AreaLevel = "province"
	AreaCity     AreaLevel = "city"
	AreaDistrict AreaLevel = "district"
)

type Area struct {
	CityID       model.LongID `json:"cityID"`
	ParentCityID model.LongID `json:"parentCityID"`
	CityCode     string       `json:"cityCode"`
	Name         string       `json:"name"`
	Py           string       `json:"py"`
	ShortPy      string       `json:"shortPy"`
	Longitude    string       `json:"longitude"`
	Latitude     string       `json:"latitude"`
	Level        AreaLevel    `json:"level"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
}

func (Area) TableName() string {
	return "business_area"
}
