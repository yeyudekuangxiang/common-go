package activity

type Zyh struct {
	Id     int64  `json:"id"`
	Openid string `json:"openid"`
	VolId  string `json:"vol_id"`
}

func (Zyh) TableName() string {
	return "activity_zyh"
}
