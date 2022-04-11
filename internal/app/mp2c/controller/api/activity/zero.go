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
	user := util.GetAuthUser(ctx)

	loginUrl, err := activity.DefaultZeroService.AutoLogin(user.ID)
	return gin.H{
		"loginUrl": loginUrl,
	}, err
}
