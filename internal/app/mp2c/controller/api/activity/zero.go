package activity

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/service/activity"
	"mio/internal/pkg/util"
)

var DefaultZeroController = ZeroController{}

type ZeroController struct {
}

func (ctr ZeroController) AutoLogin(ctx *gin.Context) (gin.H, error) {
	form := ZeroAutoLoginForm{}
	if err := util.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	user := util.GetAuthUser(ctx)

	loginUrl, err := activity.DefaultZeroService.AutoLogin(user.ID, form.Short)
	return gin.H{
		"loginUrl": loginUrl,
	}, err
}
func (ctr ZeroController) StoreUrl(ctx *gin.Context) (gin.H, error) {
	form := ZeroStoreUrlForm{}
	if err := util.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	short, err := activity.DefaultZeroService.StoreUrl(form.Url)
	return gin.H{
		"short": short,
	}, err
}
