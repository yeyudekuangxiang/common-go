package wxapp

import "time"

type AccessTokenCenter interface {
	AccessToken() (token string, expireIn time.Time, err error)
	RefreshToken(oldToken string) (token string, expireIn time.Time, err error)
	IsExpired(code string) (bool, error)
}
