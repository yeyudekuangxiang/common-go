package business

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/service/business"
	"mio/internal/pkg/util/apiutil"
	"time"
)

var DefaultCarbonController = CarbonController{}

type CarbonController struct{}

func (CarbonController) CollectEvCar(ctx *gin.Context) (gin.H, error) {
	form := CarbonCollectEvCarForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthBusinessUser(ctx)
	err := business.DefaultCarbonService.CarbonCreditEvCar(user.ID, form.Electricity)
	return nil, err
}
func (CarbonController) CollectOnlineMeeting(ctx *gin.Context) (gin.H, error) {
	form := CarbonCollectOnlineMeetingForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	duration := time.Duration(float64(time.Hour) * form.OnlineDuration)

	start := time.Now()

	user := apiutil.GetAuthBusinessUser(ctx)
	err := business.DefaultCarbonService.CarbonCreditOnlineMeeting(user.ID, duration, start, start.Add(duration))
	return nil, err
}
func (CarbonController) CollectSaveWaterElectricity(ctx *gin.Context) (gin.H, error) {
	form := CarbonCollectSaveWaterElectricityForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthBusinessUser(ctx)
	err := business.DefaultCarbonService.CarbonCreditSaveWaterElectricity(user.ID, form.Water, form.Electricity)
	return nil, err
}
func (CarbonController) CollectPublicTransport(ctx *gin.Context) (gin.H, error) {
	form := CarbonCollectPublicTransportForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthBusinessUser(ctx)
	err := business.DefaultCarbonService.CarbonCreditPublicTransport(user.ID, form.Bus, form.Metro)
	return nil, err
}
