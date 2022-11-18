package common

type CityListRequest struct {
	//CityCode    string `json:"cityCode"`
	CityPidCode string `json:"cityPidCode" form:"cityPidCode"`
}
