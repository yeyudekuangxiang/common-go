package common

type CreateCityParams struct {
	CityCode string
	Name     string
	PidCode  string
}

type GetByCityCodeParams struct {
	CityCode string
}
