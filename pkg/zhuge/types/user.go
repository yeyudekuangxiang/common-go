package types

type User struct {
	Ct   int64  `json:"$ct"`
	Cuid string `json:"$cuid"`
}

func (User) F3127563263F0A80C8007E338109E07F() {

}
func (u User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"$ct":   u.Ct,
		"$cuid": u.Cuid,
	}
}
