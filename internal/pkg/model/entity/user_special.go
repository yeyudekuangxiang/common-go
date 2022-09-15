package entity

type UserSpecial struct {
	ID       int64  `gorm:"primary_key;column:id" json:"id"`
	Phone    string `json:"phone"`
	Status   int    `json:"status"` //是否绑定 0否1是
	Identity string `json:"identity"`
	Icon     string `json:"icon"`
	Partner  int    `json:"partner"`
}

func (UserSpecial) TableName() string {
	return "user_special"
}
