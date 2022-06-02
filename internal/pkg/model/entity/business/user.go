package business

import "mio/internal/pkg/model"

type User struct {
	ID            int64      `json:"-" gorm:"primaryKey;not null;type:serial8;comment:企业用户表"`
	Uid           string     `json:"uid" gorm:"not null;type:varchar(100);comment:uuid"`
	BDepartmentId int        `json:"bDepartmentId" gorm:"not null;type:int4;comment:所属企业版部门id"`
	BCompanyId    int        `json:"bCompanyId" gorm:"not null;type:int4;comment:所属企业版企业表主键id"`
	Nickname      string     `json:"nickname" gorm:"not null;type:varchar(100);comment:昵称"`
	Mobile        string     `json:"mobile" gorm:"not null;type:varchar(20);comment:手机号"`
	TelephoneCode string     `json:"telephoneCode" gorm:"not null;type:varchar(10);default:'86';comment:国际区号 默认86"`
	Realname      string     `json:"realname" gorm:"not null;type:varchar(20);default:'';comment:真实姓名"`
	Status        int8       `json:"status" gorm:"not null;type:int2;default:1;comment:在职状态 1在职 2离职"`
	CreatedAt     model.Time `json:"createdAt" gorm:"not null;type:timestamptz"`
	UpdatedAt     model.Time `json:"updatedAt" gorm:"not null;type:timestamptz"`
}

func (User) TableName() string {
	return "business_user"
}
