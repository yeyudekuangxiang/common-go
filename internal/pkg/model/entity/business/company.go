package business

import "mio/internal/pkg/model"

type Company struct {
	ID        int        `json:"-" gorm:"primaryKey;not null;type:serial4;comment:企业表"`
	Cid       string     `json:"cid" gorm:"not null;type:varchar(100);comment:企业uuid"`
	Name      string     `json:"name" gorm:"not null;type:varchar(100);comment:企业名称"`
	Email     string     `json:"email" gorm:"not null;type:varchar(100);comment:企业邮箱"`
	Password  string     `json:"password" gorm:"not null;type:varchar(100);comment:企业密码 sha1加密"`
	CreatedAt model.Time `json:"createdAt" gorm:"not null;type:timestamp"`
	UpdatedAt model.Time `json:"updatedAt" gorm:"not null;type:timestamp"`
}

func (Company) TableName() string {
	return "business_company"
}
