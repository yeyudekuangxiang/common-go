package business

type CityLisDTO struct {
	LikeName    string
	LikePy      string
	LikeShortPy string
}

type AreaListDTO struct {
	CityCodes []string
}

type ShortArea struct {
	ID        int64  `json:"id"`
	CityCode  string `json:"cityCode"`
	Name      string `json:"name"`
	Py        string `json:"py"`
	ShortPy   string `json:"shortPy"`
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}
type CityProvince struct {
	Province ShortArea `json:"province"`
	City     ShortArea `json:"city"`
}
