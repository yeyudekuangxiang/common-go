package model

type OrderAddressInfo struct {
	Province string `json:"province" form:"province" alias:"province"`
	City     string `json:"city" form:"city" alias:"city"`
	Area     string `json:"area" form:"area" alias:"area"`
	Street   string `json:"street" form:"street" alias:"street"`
	Address  string `json:"address" form:"address" alias:"address"`
	Mobile   string `json:"mobile" form:"mobile" alias:"mobile"`
	Name     string `json:"name" form:"name" alias:"name"`
}
