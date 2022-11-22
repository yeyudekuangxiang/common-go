package wxoa

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net/http"
	"time"
)

type AccessTokenServer struct {
	Redis  *redis.Client
	AppId  string
	Secret string
	Prefix string
}

func (srv *AccessTokenServer) cacheKey() string {
	return fmt.Sprintf(srv.Prefix, srv.Secret)
}
func (srv *AccessTokenServer) getToken() (*AccessToken, error) {
	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + srv.AppId +
		"&secret=" + srv.Secret

	httpResp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http.Status: %s", httpResp.Status)
	}
	oaResponse := AccessTokenResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&oaResponse)
	if err != nil {
		return nil, err
	}
	if oaResponse.ErrCode == 0 {
		return &oaResponse.AccessToken, nil
	}
	return nil, fmt.Errorf("%d %s", oaResponse.ErrCode, oaResponse.ErrMsg)
}
func (srv *AccessTokenServer) Token() (token string, err error) {
	//return srv.RefreshToken("")
	token, err = srv.Redis.Get(context.Background(), srv.cacheKey()).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	if token != "" {
		return
	}

	accessToken, err := srv.getToken()

	if err != nil {
		return
	}
	token = accessToken.AccessToken
	expiresIn := time.Duration(accessToken.ExpiresIn-60) * time.Second
	srv.Redis.Set(context.Background(), srv.cacheKey(), token, expiresIn)
	return
}
func (srv *AccessTokenServer) RefreshToken(currentToken string) (token string, err error) {
	token, err = srv.Redis.Get(context.Background(), srv.cacheKey()).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	if token != "" && currentToken != "" && token != currentToken {
		return
	}

	accessToken, err := srv.getToken()
	if err != nil {
		return "", err
	}
	token = accessToken.AccessToken
	expiresIn := time.Duration(accessToken.ExpiresIn-60) * time.Second
	srv.Redis.Set(context.Background(), srv.cacheKey(), token, expiresIn)
	return
}
func (srv *AccessTokenServer) IID01332E16DF5011E5A9D5A4DB30FED8E1() {

}
