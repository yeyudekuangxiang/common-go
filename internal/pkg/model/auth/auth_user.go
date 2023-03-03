package auth

type User struct {
	ID     int64  `json:"id"`
	Mobile string `json:"mobile"`
}

func (au User) Valid() error {
	return nil
}
