package points

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/service/point"
	"mio/internal/pkg/util/apiutil"
)

var PointsCollectController = PointCollectController{}

type PointCollectController struct {
}

func (PointCollectController) Collect(ctx *gin.Context) (gin.H, error) {
	form := api.PointCollectForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)
	client := point.NewClientHandle(context.NewMioContext(), user.OpenId, form.ImgUrl)
	if err := client.HandleCollectCommand(form.PointCollectType); err != nil {
		return nil, err
	}
	return nil, nil
}
