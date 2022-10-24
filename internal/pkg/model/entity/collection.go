package entity

import "time"

type Collection struct {
	Id        int64     `gorm:"primary_key;column:id" json:"id"`
	ObjId     int64     `gorm:"column:obj_id" json:"objId"`
	ObjType   int       `gorm:"column:obj_type" json:"objType"`
	Status    int       `gorm:"column:status" json:"status"`
	OpenId    string    `gorm:"column:open_id" json:"openId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (c Collection) TableName() string {
	return "collection"
}
