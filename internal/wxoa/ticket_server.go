package wxoa

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"mio/config"
	"mio/core/app"
	"net/http"
	"time"
)

type TicketTokenServer struct {
	Redis       *redis.Client
	AppId       string
	TokenServer *AccessTokenServer
}

func (srv *TicketTokenServer) cacheKey() string {
	return fmt.Sprintf(config.RedisKey.AccessToken, "oaticket", srv.AppId)
}
func (srv *TicketTokenServer) getTicket() (*Ticket, error) {
	accessToken, err := srv.TokenServer.Token()
	if err != nil {
		return nil, err
	}
	url := "https://api.weixin.qq.com/cgi-bin/ticket/getticket?type=jsapi&access_token=" + accessToken

	httpResp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http.Status: %s", httpResp.Status)
	}
	oaResponse := TicketResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&oaResponse)
	if err != nil {
		return nil, err
	}
	if oaResponse.ErrCode == 0 {
		return &oaResponse.Ticket, nil
	}
	return nil, fmt.Errorf("%d %s", oaResponse.ErrCode, oaResponse.ErrMsg)
}
func (srv *TicketTokenServer) Ticket() (ticket string, err error) {
	ticket, err = srv.Redis.Get(context.Background(), srv.cacheKey()).Result()
	if err != nil {
		app.Logger.Error(err)
	}
	if ticket != "" {
		return
	}

	ticketData, err := srv.getTicket()

	if err != nil {
		return
	}
	ticket = ticketData.Ticket
	expiresIn := time.Duration(ticketData.ExpiresIn-60) * time.Second
	srv.Redis.Set(context.Background(), srv.cacheKey(), ticket, expiresIn)
	return
}
func (srv *TicketTokenServer) RefreshTicket(currentTicket string) (ticket string, err error) {
	ticket, err = srv.Redis.Get(context.Background(), srv.cacheKey()).Result()
	if err != nil {
		app.Logger.Error(err)
	}
	if ticket != "" && ticket != currentTicket {
		return
	}

	ticketData, err := srv.getTicket()
	if err != nil {
		return
	}
	ticket = ticketData.Ticket
	expiresIn := time.Duration(ticketData.ExpiresIn-60) * time.Second
	srv.Redis.Set(context.Background(), srv.cacheKey(), ticket, expiresIn)
	return
}
func (srv *TicketTokenServer) IIDB04E44A0E1DC11E5ADCEA4DB30FED8E1() {

}
