package auth

import (
	"mio/internal/pkg/model"
)

type BusinessUser struct {
	ID        int64  `json:"id"`
	Mobile    string `json:"mobile"`
	Uid       string `json:"uid"`
	CreatedAt model.Time
}

func (au BusinessUser) Valid() error {
	return nil
}
