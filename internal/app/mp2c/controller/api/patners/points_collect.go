package patners

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api"
	"mio/internal/pkg/service/point"
	"mio/internal/pkg/util/apiutil"
)

var PatnersPointCollectController = PointCollectController{}

type PointCollectController struct {
}

func (PointCollectController) Collect(ctx *gin.Context) (gin.H, error) {
	form := api.PointCollectForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)
	client := point.NewClientHandle(user.OpenId, form.ImgUrl)
	client.HandleCommand(form.PointCollectType)
	return nil, nil
}
