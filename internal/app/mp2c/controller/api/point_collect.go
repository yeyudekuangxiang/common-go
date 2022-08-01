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
	user := apiutil.GetAuthUser(ctx)
	switch service.PointCollectType(form.PointCollectType) {
	case service.PointCollectBikeRideType:
		point, err = service.DefaultPointCollectService.CollectBikeRide(user.OpenId, form.ImgUrl)
	case service.PointCollectCoffeeCupType:
		point, err = service.DefaultPointCollectService.CollectCoffeeCup(user.OpenId, form.ImgUrl)
	case service.PointCollectPowerReplaceType:
		point, err = service.DefaultPointCollectService.CollectPowerReplace(user.OpenId, form.ImgUrl)
	}
	return gin.H{
		"point": point,
	}, err
}
