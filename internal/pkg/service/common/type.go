package common

type CreateCityParams struct {
	CityCode string
	Name     string
	PidCode  string
}

type GetCityParams struct {
	CityCode string
	CityName string
}

type GetCityListParams struct {
	//CityCode    string
	CityPidCode string
}
