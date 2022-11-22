package types

type Revenue struct {
	Ct              int64   `json:"$ct"`
	Eid             string  `json:"$eid"`
	Cuid            string  `json:"$cuid"`
	Sid             int64   `json:"$sid"`
	Vn              string  `json:"$vn"`
	Cn              string  `json:"$cn"`
	Cr              int     `json:"$cr"`
	Os              string  `json:"$os"`
	Net             int     `json:"$net"`
	Price           float64 `json:"$price"`
	ProductID       string  `json:"$productID"`
	ProductQuantity int     `json:"$productQuantity"`
	RevenueType     string  `json:"$revenueType"`
}

func (Revenue) F3127563263F0A80C8007E338109E07F() {

}
func (r Revenue) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"$ct":              r.Ct,
		"$eid":             r.Eid,
		"$cuid":            r.Cuid,
		"$sid":             r.Sid,
		"$vn":              r.Vn,
		"$cn":              r.Cn,
		"$cr":              r.Cr,
		"$os":              r.Os,
		"$net":             r.Net,
		"$price":           r.Price,
		"$productID":       r.ProductID,
		"$productQuantity": r.ProductQuantity,
		"$revenueType":     r.RevenueType,
	}
}
