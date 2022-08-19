package business

import (
	"mio/internal/pkg/model"
	ebusiness "mio/internal/pkg/model/entity/business"
)

type CityProvinceListDTO struct {
	Search string
}

type AreaListDTO struct {
	Search      string //根据 name、py、ShortPy 模糊搜索
	CityIds     []int64
	CityCodes   []string
	LikeName    string
	LikePy      string
	LikeShortPy string
	Level       ebusiness.AreaLevel
	Names       []string
}

type ShortArea struct {
	CityID    model.LongID `json:"cityID"`
	CityCode  string       `json:"cityCode"`
	Name      string       `json:"name"`
	Py        string       `json:"py"`
	ShortPy   string       `json:"shortPy"`
	Longitude string       `json:"longitude"`
	Latitude  string       `json:"latitude"`
}
type CityProvince struct {
	Province ShortArea `json:"province"`
	City     ShortArea `json:"city"`
}
type GroupCityProvince struct {
	Letter string         `json:"letter"`
	Items  []CityProvince `json:"items"`
}
