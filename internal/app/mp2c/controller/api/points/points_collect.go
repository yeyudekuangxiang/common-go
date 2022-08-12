package points

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/service/point"
	"mio/internal/pkg/util/apiutil"
)

var DefaultPointsCollectController = PointCollectController{}

type PointCollectController struct {
}

func (ctr PointCollectController) Collect(ctx *gin.Context) (gin.H, error) {
	form := api.PointCollectForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)
	client := point.NewClientHandle(context.NewMioContext(), &point.ClientHandle{
		OpenId: user.OpenId,
		ImgUrl: form.ImgUrl,
		Type:   point.CollectType(form.PointCollectType),
	})
	if err := client.HandleCollectCommand(form.PointCollectType); err != nil {
		return nil, err
	}
	return nil, nil
}

func (ctr PointCollectController) CallCollect(ctx *gin.Context) (gin.H, error) {
	form := api.PointCollectForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)
	client := point.NewClientHandle(context.NewMioContext(), &point.ClientHandle{
		OpenId: user.OpenId,
		ImgUrl: form.ImgUrl,
		Type:   point.CollectType(form.PointCollectType),
	})
	if err := client.HandleCollectCommand(form.PointCollectType); err != nil {
		return nil, err
	}
	return nil, nil
}

func (ctr PointCollectController) GetPageData(ctx *gin.Context) (gin.H, error) {
	form := api.CollectType{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)
	client := point.NewClientHandle(context.NewMioContext(), &point.ClientHandle{
		OpenId: user.OpenId,
		Type:   point.CollectType(form.PointCollectType),
	})
	res, err := client.HandlePageDataCommand()
	if err != nil {
		return nil, err
	}
	return gin.H{
		"data": res,
	}, nil
}
