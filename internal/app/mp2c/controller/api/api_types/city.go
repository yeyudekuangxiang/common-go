package api_types

// DTO

type CreateCityDto struct {
	CityCode string
	Name     string
	PidCode  string
}

type GetByCityCode struct {
	CityCode string
}
