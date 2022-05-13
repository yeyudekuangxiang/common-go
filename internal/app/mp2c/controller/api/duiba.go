package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
	duibaApi "mio/pkg/duiba/api/model"
)

var DefaultDuiBaController = DuiBaController{}

type DuiBaController struct {
}

func (DuiBaController) AutoLogin(ctx *gin.Context) (gin.H, error) {
	form := DuibaAutoLoginForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)
	u, err := service.DefaultDuiBaService.AutoLogin(service.AutoLoginParam{
		UserId: user.ID,
		Path:   form.Path,
	})
	if err != nil {
		return nil, err
	}
	return gin.H{
		"loginUrl": u,
	}, nil
}

func (DuiBaController) ExchangeCallback(ctx *gin.Context) gin.H {
	form := duibaApi.Exchange{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return gin.H{
			"status":       "fail",
			"errorMessage": err.Error(),
			"credits":      "0",
		}
	}

	err := service.DefaultDuiBaService.CheckSign(form)
	if err != nil {
		app.Logger.Error("ExchangeCallback 参数验证失败", form, err)
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
	form := duibaApi.ExchangeResult{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		app.Logger.Error("ExchangeResultNoticeCallback 参数获取失败", ctx, err)
		return err.Error()
	}
	err := service.DefaultDuiBaService.CheckSign(form)
	if err != nil {
		app.Logger.Error("ExchangeResultNoticeCallback 参数验证失败", form, err)
		return err.Error()
	}
	err = service.DefaultDuiBaService.ExchangeResultNoticeCallback(form)
	if err != nil {
		app.Logger.Error("ExchangeResultNoticeCallback 退还积分失败", form, err)
		return err.Error()
	}
	return "ok"
}
func (DuiBaController) OrderCallback(ctx *gin.Context) string {

	form := duibaApi.OrderInfo{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		app.Logger.Error("OrderCallback 参数获取失败", ctx, err)
		return err.Error()
	}
	err := service.DefaultDuiBaService.CheckSign(form)
	if err != nil {
		app.Logger.Error("OrderCallback 参数验证失败", form, err)
		return err.Error()
	}
	err = service.DefaultDuiBaService.OrderCallback(form)
	if err != nil {
		app.Logger.Error("OrderCallback 同步订单失败", form, err)
		return err.Error()
	}
	return "ok"
}
func (DuiBaController) PointAddLogCallback(ctx *gin.Context) gin.H {
	/*d, err := ioutil.ReadAll(ctx.Request.Body)
	fmt.Println(string(d), err)*/

	form := duibaApi.PointAdd{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return gin.H{
			"status":       "fail",
			"errorMessage": err.Error(),
			"credits":      0,
		}
	}

	userPoint, _ := service.DefaultPointService.FindByOpenId(form.Uid)
	err := service.DefaultDuiBaService.CheckSign(form)
	if err != nil {
		return gin.H{
			"status":       "fail",
			"errorMessage": err.Error(),
			"credits":      userPoint.Balance,
		}
	}

	tranId, err := service.DefaultDuiBaService.PointAddCallback(form)
	userPoint, _ = service.DefaultPointService.FindByOpenId(form.Uid)
	if err != nil {
		return gin.H{
			"status":       "fail",
			"errorMessage": err.Error(),
			"credits":      userPoint.Balance,
		}
	}

	return gin.H{
		"status":       "ok",
		"errorMessage": "",
		"bizId":        tranId,
		"credits":      userPoint.Balance,
	}
}
