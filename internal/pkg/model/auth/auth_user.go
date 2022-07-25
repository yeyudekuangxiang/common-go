package auth

import (
	"mio/internal/pkg/model"
)

type User struct {
	ID        int64      `json:"id"`
	Mobile    string     `json:"mobile"`
	CreatedAt model.Time `json:"createdAt"`
}

func (au User) Valid() error {
	return nil
}
