package repotypes

type GetCityListDO struct {
	CityCode      string
	CityCodeSlice []string
	Name          string
	PidCode       string
}

type GetCityByCode struct {
	CityCode string
}
