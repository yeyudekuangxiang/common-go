package entity

type UserSpecial struct {
	ID     int64  `gorm:"primary_key;column:id" json:"id"`
	Phone  string `json:"phone"`
	Status int    `json:"status"` //是否绑定 0否1是
}

func (UserSpecial) TableName() string {
	return "user_special"
}
