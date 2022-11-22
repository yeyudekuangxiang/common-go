package types

type Device struct {
	Ct   int64  `json:"$ct"`
	Br   string `json:"$br"`
	Dv   string `json:"$dv"`
	Imei string `json:"$imei"`
	Rs   string `json:"$rs"`
}

func (Device) F3127563263F0A80C8007E338109E07F() {

}
func (d Device) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"$ct":   d.Ct,
		"$br":   d.Br,
		"$dv":   d.Dv,
		"$imei": d.Imei,
		"$rs":   d.Rs,
	}
}
