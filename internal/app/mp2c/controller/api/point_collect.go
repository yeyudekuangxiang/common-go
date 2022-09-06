package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
)

var DefaultPointCollectController = PointCollectController{}

type PointCollectController struct {
}

func (PointCollectController) Collect(ctx *gin.Context) (gin.H, error) {
	form := PointCollectForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	var (
		point int
		err   error
	)
	user := apiutil.GetAuthUser(ctx)
	switch service.PointCollectType(form.PointCollectType) {
	case service.PointCollectBikeRideType:
		point, err = service.DefaultPointCollectService.CollectBikeRide(user.OpenId, user.Risk, form.ImgUrl)
	case service.PointCollectCoffeeCupType:
		point, err = service.DefaultPointCollectService.CollectCoffeeCup(user.OpenId, user.Risk, form.ImgUrl)
	case service.PointCollectPowerReplaceType:
		point, err = service.DefaultPointCollectService.CollectPowerReplace(user.OpenId, user.Risk, form.ImgUrl)
	case service.PointCollectReducePlastic:
		point, err = service.DefaultPointCollectService.CollectReducePlastic(user.OpenId, user.Risk, form.ImgUrl)
	}
	carbon := 0.0
	if err == nil {
		//发碳量
		carbon, err = service.NewCarbonTransactionService(context.NewMioContext()).Create(api_types.CreateCarbonTransactionDto{
			OpenId:  user.OpenId,
			UserId:  user.ID,
			Type:    entity.CarbonTransactionType(form.PointCollectType),
			Value:   1,
			Info:    fmt.Sprintf("{imageUrl=%s}", form.ImgUrl),
			AdminId: 0,
			Ip:      ctx.ClientIP(),
		})
	}
	return gin.H{
		"point":  point,
		"carbon": carbon,
	}, err
}
