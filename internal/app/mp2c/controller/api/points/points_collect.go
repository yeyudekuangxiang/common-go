package points

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service/point"
	"mio/internal/pkg/util/apiutil"
)

var DefaultPointsCollectController = PointCollectController{}

type PointCollectController struct {
}

// ImageCollect 根据图片收集积分
func (ctr PointCollectController) ImageCollect(ctx *gin.Context) (gin.H, error) {
	form := api.PointCollectForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)
	client := point.NewClientHandle(context.NewMioContext(), &point.ClientHandle{
		OpenId: user.OpenId,
		ImgUrl: form.ImgUrl,
		Type:   entity.PointTransactionType(form.PointCollectType),
	})
	result, err := client.HandleImageCollectCommand()
	if err != nil {
		return nil, err
	}
	return gin.H{
		"data": result,
	}, nil
}

// CallCollect 收集积分结束后返回数据
func (ctr PointCollectController) CallCollect(ctx *gin.Context) (gin.H, error) {
	form := api.PointCollectForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)
	client := point.NewClientHandle(context.NewMioContext(), &point.ClientHandle{
		OpenId: user.OpenId,
		ImgUrl: form.ImgUrl,
		Type:   entity.PointTransactionType(form.PointCollectType),
	})
	result, err := client.HandlePageDataCommand()
	if err != nil {
		return nil, err
	}
	return gin.H{"data": result}, nil
}

// GetPageData 收集积分前返回数据
func (ctr PointCollectController) GetPageData(ctx *gin.Context) (gin.H, error) {
	form := api.CollectType{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)
	client := point.NewClientHandle(context.NewMioContext(), &point.ClientHandle{
		OpenId: user.OpenId,
		Type:   entity.PointTransactionType(form.PointCollectType),
	})

	res, err := client.HandlePageDataCommand()
	if err != nil {
		return nil, err
	}
	return gin.H{
		"data": res,
	}, nil
}
