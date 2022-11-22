package authtool

import (
	"github.com/dgrijalva/jwt-go"
)

func ParseJwtToken(tokenKey string, token string, v jwt.Claims) error {
	_, err := jwt.ParseWithClaims(token, v, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenKey), nil
	})
	return err
}

func CreateJwtToken(tokenKey string, v jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, v)
	return token.SignedString([]byte(tokenKey))
}
