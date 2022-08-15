package business

import "time"

type Area struct {
	ID        int64     `json:"id"`
	CityCode  string    `json:"cityCode"`
	Name      string    `json:"name"`
	PidCode   string    `json:"pidCode"`
	Py        string    `json:"py"`
	ShortPy   string    `json:"shortPy"`
	Longitude string    `json:"longitude"`
	Latitude  string    `json:"latitude"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
