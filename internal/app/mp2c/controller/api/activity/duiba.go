package activity

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/medivhzhan/weapp/v3"
	"io/ioutil"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/activity"
	"mio/internal/pkg/util/apiutil"
	"strconv"
)

var DefaultZeroController = ZeroController{}

type ZeroController struct {
}

func (ctr ZeroController) AutoLogin(ctx *gin.Context) (gin.H, error) {
	form := ZeroAutoLoginForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthUser(ctx)

	loginUrl, err := activity.DefaultZeroService.AutoLogin(user.ID, form.Short)
	return gin.H{
		"loginUrl": loginUrl,
	}, err
}
func (ctr ZeroController) StoreUrl(ctx *gin.Context) (gin.H, error) {
	form := ZeroStoreUrlForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	short, err := activity.DefaultZeroService.StoreUrl(form.Url)
	return gin.H{
		"short": short,
	}, err
}

func (ctr ZeroController) DuiBaAutoLogin(ctx *gin.Context) (gin.H, error) {
	form := DuiBaAutoLoginForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthUser(ctx)

	loginUrl, err := activity.DefaultZeroService.DuiBaAutoLogin(user.ID, form.ActivityId, form.Short, form.ThirdParty, ctx.ClientIP())

	return gin.H{
		"loginUrl": loginUrl,
	}, err
}

func (ctr ZeroController) DuiBaStoreUrl(ctx *gin.Context) (gin.H, error) {
	form := DuiBaStoreUrlForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	short, err := activity.DefaultZeroService.DuiBaStoreUrl(form.ActivityId, form.Url)
	return gin.H{
		"short": short,
	}, err
}
func (ctr ZeroController) GetActivityMiniQR(ctx *gin.Context) error {
	form := GetDuiBaActivityQrForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return err
	}
	if form.Password != "greencatniubi123..." {
		return errors.New("密码错误")
	}

	activity, err := service.DefaultDuiBaActivityService.FindActivity(form.ActivityId)
	if err != nil {
		return err
	}

	res, err1, err2 := app.Weapp.GetQRCode(&weapp.QRCode{
		Path: fmt.Sprintf("/pages/duiba_v2/duiba/index?activityId=%s", activity.ActivityId),
	})
	if err2 != nil {
		return err2
	}

	if err1.ErrCode != 0 {
		return errors.New(strconv.Itoa(err1.ErrCode) + err1.ErrMSG)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	ctx.Writer.Write(body)
	ctx.Header(res.Header.Get("content-type"), "image/jpeg")
	ctx.Abort()
	return nil
}
