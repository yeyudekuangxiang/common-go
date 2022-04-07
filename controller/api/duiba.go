package api

import (
	"github.com/gin-gonic/gin"
	"mio/core/app"
	"mio/internal/util"
	duibaApi "mio/pkg/duiba/api"
	"mio/service"
)

var DefaultDuiBaController = DuiBaController{}

type DuiBaController struct {
}

func (DuiBaController) AutoLogin(ctx *gin.Context) (gin.H, error) {
	form := DuibaAutoLoginForm{}
	if err := util.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := util.GetAuthUser(ctx)
	u, err := service.DefaultDuiBaService.AutoLogin(user.ID, form.Path)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"loginUrl": u,
	}, nil
}

func (DuiBaController) ExchangeCallback(ctx *gin.Context) gin.H {
	form := duibaApi.ExchangeForm{}
	if err := util.BindForm(ctx, &form); err != nil {
		return gin.H{
			"status":       "fail",
			"errorMessage": err.Error(),
			"credits":      "0",
		}
	}
	result, err := service.DefaultDuiBaService.ExchangeCallback(form)
	if err != nil {
		return gin.H{
			"status":       "fail",
			"errorMessage": err.Error(),
			"credits":      "0",
		}
	}
	return gin.H{
		"status":       "ok",
		"errorMessage": " ",
		"bizId":        result.BizId,
		"credits":      result.Credits,
	}
}

func (DuiBaController) ExchangeResultNoticeCallback(ctx *gin.Context) string {
	form := duibaApi.ExchangeResultForm{}
	if err := util.BindForm(ctx, &form); err != nil {
		app.Logger.Error("ExchangeResultNoticeCallback 参数获取失败", ctx, err)
		return "ok"
	}
	err := service.DefaultDuiBaService.ExchangeResultNoticeCallback(form)
	if err != nil {
		app.Logger.Error("ExchangeResultNoticeCallback 退还积分失败", form, err)
	}
	return "ok"
}
