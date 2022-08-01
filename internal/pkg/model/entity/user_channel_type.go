package entity

type UserChannelType struct {
	ID          int64  `gorm:"primary_key;column:id" json:"id"`
	Name        string `gorm:"column:name" json:"name"`
	Description string `gorm:"column:description" json:"description"`
}

func (UserChannelType) TableName() string {
	return "user_channel_type"
}
