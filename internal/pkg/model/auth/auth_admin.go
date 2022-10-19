package auth

import (
	"mio/internal/pkg/model"
)

type Admin struct {
	MioAdminID int `json:"mioAdminId"`
}

func (au Admin) Valid() error {
	return nil
}

type OldAdmin struct {
	UserId model.StrToInt `json:"userId"`
}
