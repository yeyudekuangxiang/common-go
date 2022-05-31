package business

import "mio/internal/pkg/model"

type Department struct {
	ID         int        `json:"id" gorm:"primaryKey;not null;type:serial4;comment:部门表"`
	Title      string     `json:"title" gorm:"not null;type:varchar(20);comment:部门标题"`
	BCompanyId int        `json:"bCompanyId" gorm:"not null;type:int4;comment:部门所属企业id"`
	Icon       string     `json:"icon" gorm:"not null;type:varchar(500);comment:部门图标"`
	CreatedAt  model.Time `json:"createdAt" gorm:"not null;type:timestamptz"`
	UpdatedAt  model.Time `json:"updatedAt" gorm:"not null;type:timestamptz"`
}

func (Department) TableName() string {
	return "business_department"
}
