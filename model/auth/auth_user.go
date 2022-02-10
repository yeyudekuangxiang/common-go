package auth

type User struct {
	BindMobile string `json:"bind_mobile"`
	CreatedAt  string `json:"created_at"`
	Guid       string `json:"guid"`
	NowOpenid  string `json:"now_openid"`
	Session    string `json:"session"`
	SessionKey string `json:"session_key"`
}

func (au User) Valid() error {
	return nil
}
