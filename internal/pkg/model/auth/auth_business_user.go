package auth

import "time"

type BusinessUser struct {
	Uid       string `json:"uid"`
	CreatedAt time.Time
}

func (au BusinessUser) Valid() error {
	return nil
}
