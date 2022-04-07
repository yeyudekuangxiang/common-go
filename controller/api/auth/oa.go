package auth

import (
	"github.com/chanxuehong/wechat/mp/jssdk"
	"github.com/gin-gonic/gin"
	"mio/config"
	"mio/core/app"
	"mio/internal/errno"
	"mio/internal/util"
	"mio/internal/wxoa"
	"strconv"
	"time"
)

var DefaultOaController = OaController{}

type OaController struct {
}

func (OaController) Sign(c *gin.Context) (gin.H, error) {
	form := ConfigSignForm{}
	if err := util.BindForm(c, &form); err != nil {
		return nil, err
	}

	tickerServer := wxoa.TicketTokenServer{
		TokenServer: &wxoa.AccessTokenServer{
			AppId:  config.Config.Wxoa.AppId,
			Secret: config.Config.Wxoa.Secret,
			Redis:  app.Redis,
		},
		AppId: config.Config.Wxoa.AppId,
		Redis: app.Redis,
	}

	ticker, err := tickerServer.Ticket()
	if err != nil {
		app.Logger.Error(err)
		return nil, errno.InternalServerError
	}

	nonceStr := util.Md5(time.Now().String())
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sign := jssdk.WXConfigSign(ticker, nonceStr, timestamp, form.Url)
	return gin.H{
		"appId":     config.Config.Wxoa.AppId,
		"timestamp": timestamp,
		"nonceStr":  nonceStr,
		"signature": sign,
	}, nil
}
