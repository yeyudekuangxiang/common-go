package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.miotech.com/miotech-application/backend/common-go/duiba"
	duibaApi "gitlab.miotech.com/miotech-application/backend/common-go/duiba/api/model"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/service"
	duiba2 "mio/internal/pkg/service/duiba"
	"mio/internal/pkg/util/apiutil"
	"net/http"
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
	u, err := duiba2.DefaultDuiBaService.AutoLogin(service.AutoLoginParam{
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

	err := duiba2.DefaultDuiBaService.CheckSign(form)
	if err != nil {
		app.Logger.Error("ExchangeCallback 参数验证失败", form, err)
		return gin.H{
			"status":       "fail",
			"errorMessage": err.Error(),
			"credits":      "0",
		}
	}

	result, err := duiba2.DefaultDuiBaService.ExchangeCallback(form)
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
	err := duiba2.DefaultDuiBaService.CheckSign(form)
	if err != nil {
		app.Logger.Error("ExchangeResultNoticeCallback 参数验证失败", form, err)
		return err.Error()
	}
	err = duiba2.DefaultDuiBaService.ExchangeResultNoticeCallback(form)
	if err != nil {
		app.Logger.Error("ExchangeResultNoticeCallback 退还积分失败", form, err)
		return err.Error()
	}
	return "ok"
}
func (DuiBaController) VirtualGoodCallback(ctx *gin.Context) gin.H {
	form := duibaApi.VirtualGood{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		app.Logger.Error("VirtualGoodCallback 参数获取失败", ctx, err)
		return gin.H{
			"status":        "fail",
			"credits":       0,
			"supplierBizId": "",
		}
	}
	err := duiba2.DefaultDuiBaService.CheckSign(form)
	if err != nil {
		app.Logger.Error("VirtualGoodCallback 参数验证失败", form, err)
		return gin.H{
			"status":        "fail",
			"credits":       0,
			"supplierBizId": "",
		}
	}
	orderId, credit, err := duiba2.DefaultDuiBaService.VirtualGoodCallback(form)
	if err != nil {
		app.Logger.Error("VirtualGoodCallback 兑换虚拟商品失败", form, err)
		return gin.H{
			"status":        "fail",
			"credits":       0,
			"supplierBizId": "",
		}
	}
	return gin.H{
		"status":        "success",
		"credits":       credit,
		"supplierBizId": orderId,
	}
}
func (DuiBaController) OrderCallback(ctx *gin.Context) string {

	form := duibaApi.OrderInfo{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		app.Logger.Error("OrderCallback 参数获取失败", ctx, err)
		return err.Error()
	}
	err := duiba2.DefaultDuiBaService.CheckSign(form)
	if err != nil {
		app.Logger.Error("OrderCallback 参数验证失败", form, err)
		return err.Error()
	}

	err = duiba2.DefaultDuiBaService.OrderCallback(form)
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

	pointService := service.NewPointService(context.NewMioContext())
	userPoint, _ := pointService.FindByOpenId(form.Uid)
	err := duiba2.DefaultDuiBaService.CheckSign(form)
	if err != nil {
		return gin.H{
			"status":       "fail",
			"errorMessage": err.Error(),
			"credits":      userPoint.Balance,
		}
	}

	tranId, err := duiba2.DefaultDuiBaService.PointAddCallback(form)

	userPoint, _ = pointService.FindByOpenId(form.Uid)
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
func (DuiBaController) DuiBaNoLoginH5(ctx *gin.Context) {
	form := DuiBaNoLoginH5Form{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		ctx.Status(404)
		return
	}
	duiBaService := service.NewDuiBaActivityService(context.NewMioContext())
	activity, err := duiBaService.FindActivity(form.ActivityId)
	if err != nil {
		app.Logger.Error("DuiBaNoLoginH5", form, err)
		ctx.Status(404)
		return
	}
	if activity.ID == 0 {
		ctx.Status(404)
		return
	}
	client := duiba.NewClient(config.Config.DuiBa.AppKey, config.Config.DuiBa.AppSecret)
	url, err := client.AutoLogin(duiba.AutoLoginParam{
		Uid:      "not_login",
		Redirect: activity.ActivityUrl,
	})
	if err != nil {
		app.Logger.Error("DuiBaNoLoginH5", form, err)
		ctx.Status(404)
		return
	}
	ctx.Redirect(http.StatusFound, url)
	return
}
