package business

import ebusiness "mio/internal/pkg/model/entity/business"

type AreaListPO struct {
	Search      string
	CityCodes   []string
	CityIds     []int64
	LikeName    string
	LikePy      string
	LikeShortPy string
	Level       ebusiness.AreaLevel
	Names       []string
}
type GetAreaPO struct {
	Level  ebusiness.AreaLevel
	Name   string
	CityId int64
}
