package points

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	imgCollect "mio/internal/pkg/service/point"
	"mio/internal/pkg/util/apiutil"
	"strings"
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

	uInfo := apiutil.GetAuthUser(ctx)
	client := imgCollect.NewClientHandle(
		context.NewMioContext(),
		imgCollect.ClientOptionNew(
			imgCollect.WithClientOpenId(uInfo.OpenId),
			imgCollect.WithClientUserId(uInfo.ID),
			imgCollect.WithClientImgUrl(form.ImgUrl),
			imgCollect.WithClientType(entity.PointTransactionType(strings.ToUpper(form.PointCollectType))),
		),
	)
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
	uInfo := apiutil.GetAuthUser(ctx)
	client := imgCollect.NewClientHandle(context.NewMioContext(), imgCollect.ClientOptionNew(
		imgCollect.WithClientOpenId(uInfo.OpenId),
		imgCollect.WithClientUserId(uInfo.ID),
		imgCollect.WithClientImgUrl(form.ImgUrl),
		imgCollect.WithClientType(entity.PointTransactionType(strings.ToUpper(form.PointCollectType))),
	))
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

	uInfo := apiutil.GetAuthUser(ctx)
	client := imgCollect.NewClientHandle(context.NewMioContext(), imgCollect.ClientOptionNew(
		imgCollect.WithClientOpenId(uInfo.OpenId),
		imgCollect.WithClientUserId(uInfo.ID),
		imgCollect.WithClientType(entity.PointTransactionType(strings.ToUpper(form.PointCollectType))),
	))

	res, err := client.HandlePageDataCommand()
	if err != nil {
		return nil, err
	}
	return gin.H{
		"data": res,
	}, nil
}
