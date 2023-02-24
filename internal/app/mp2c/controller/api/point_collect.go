package api

import (
	"github.com/gin-gonic/gin"
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

	//入参保存
	defer trackBehaviorInteraction(trackInteractionParam{
		Tp:   form.PointCollectType,
		Data: form,
		Ip:   ctx.ClientIP(),
	})

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

	return gin.H{
		"point":  point,
		"carbon": carbon,
	}, err
}
