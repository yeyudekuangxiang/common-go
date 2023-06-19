package auth

import (
	"mio/internal/pkg/model"
)

type Admin struct {
	Type string `json:"type"`
	ID   int64  `json:"id"`
	//兼容auth改版前获取管理员id
	MioAdminID int64 `json:"mioAdminId"`
	CreatedAt  int64 `json:"createdAt"`
}

func (au Admin) Valid() error {
	return nil
}

type OldAdmin struct {
	UserId model.StrToInt `json:"userId"`
}
