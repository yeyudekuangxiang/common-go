package wxoa

import (
	"github.com/chanxuehong/wechat/mp/core"
	"github.com/chanxuehong/wechat/mp/jssdk"
	mpoauth2 "github.com/chanxuehong/wechat/mp/oauth2"
	"github.com/chanxuehong/wechat/oauth2"
	"github.com/go-redis/redis/v8"
	"net/http"
	"sync"
	"time"
)

type WxOA struct {
	appId       string
	secret      string
	redis       *redis.Client
	cachePrefix string
}

func NewWxOA(appId string, secret string, cachePrefix string, client *redis.Client) *WxOA {
	if cachePrefix == "" {
		cachePrefix = "wxoa"
	}
	return &WxOA{appId: appId, secret: secret, cachePrefix: cachePrefix, redis: client}
}

func (oa WxOA) Oauth2() *oauth2.Client {
	return &oauth2.Client{
		Endpoint: mpoauth2.NewEndpoint(oa.appId, oa.secret),
	}
}
func (oa WxOA) MP() *MP {
	return &MP{
		appId:       oa.appId,
		secret:      oa.secret,
		redis:       oa.redis,
		cachePrefix: oa.cachePrefix + ":oa",
	}
}

type MP struct {
	appId       string
	secret      string
	redis       *redis.Client
	server      *core.Server
	mutex       sync.Mutex
	cachePrefix string
}

func (mp *MP) CoreClient() *core.Client {
	return core.NewClient(mp.AccessToken(), &http.Client{
		Timeout: 10 * time.Second,
	})
}
func (mp *MP) AccessToken() core.AccessTokenServer {
	return &AccessTokenServer{
		Redis:  mp.redis,
		AppId:  mp.appId,
		Secret: mp.secret,
		Prefix: mp.cachePrefix + ":accesstoken:",
	}
}
func (mp *MP) Ticket() jssdk.TicketServer {
	return &TicketServer{
		Redis:       mp.redis,
		AppId:       mp.appId,
		TokenServer: mp.AccessToken(),
		Prefix:      mp.cachePrefix + ":ticket:",
	}
}
