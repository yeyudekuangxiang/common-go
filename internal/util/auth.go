package util

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"mio/core/app"
)

func ParseToken(token string, v jwt.Claims) error {
	_, err := jwt.ParseWithClaims(token, v, func(token *jwt.Token) (interface{}, error) {
		tokenKey, err := GetAppConfig("app.TokenKey")
		if err != nil {
			app.Logger.Error(err)
			return "", errors.New("系统错误,请联系管理员")
		}
		return []byte(tokenKey), nil
	})
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	return nil
}

func CreateToken(v jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, v)
	tokenKey, err := GetAppConfig("app.TokenKey")
	if err != nil {
		app.Logger.Error(err)
		return "", errors.New("系统错误,请联系管理员")
	}
	return token.SignedString([]byte(tokenKey))
}
