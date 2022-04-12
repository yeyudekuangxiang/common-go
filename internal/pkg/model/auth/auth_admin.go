package auth

import (
	"mio/internal/pkg/model"
)

type Admin struct {
	ID int
}

func (au Admin) Valid() error {
	return nil
}

type OldAdmin struct {
	UserId model.StrToInt `json:"userId"`
}
