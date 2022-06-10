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
	result, err := business.DefaultCarbonService.CarbonCreditEvCar(user.ID, form.Electricity)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"carbonCredit": result.Credit,
		"point":        result.Point,
	}, err
}
func (CarbonController) CollectOnlineMeeting(ctx *gin.Context) (gin.H, error) {
	form := CarbonCollectOnlineMeetingForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	oneCityDuration := time.Duration(float64(time.Hour) * form.OneCityDuration)
	manyCityDuration := time.Duration(float64(time.Hour) * form.ManyCityDuration)

	user := apiutil.GetAuthBusinessUser(ctx)
	result, err := business.DefaultCarbonService.CarbonCreditOnlineMeeting(user.ID, oneCityDuration, manyCityDuration)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"carbonCredit": result.Credit,
		"point":        result.Point,
	}, err
}
func (CarbonController) CollectSaveWaterElectricity(ctx *gin.Context) (gin.H, error) {
	form := CarbonCollectSaveWaterElectricityForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthBusinessUser(ctx)
	result, err := business.DefaultCarbonService.CarbonCreditSaveWaterElectricity(user.ID, form.Water, form.Electricity)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"carbonCredit": result.Credit,
		"point":        result.Point,
	}, err
}
func (CarbonController) CollectPublicTransport(ctx *gin.Context) (gin.H, error) {
	form := CarbonCollectPublicTransportForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthBusinessUser(ctx)
	result, err := business.DefaultCarbonService.CarbonCreditPublicTransport(user.ID, form.Bus, form.Metro)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"carbonCredit": result.Credit,
		"point":        result.Point,
	}, err
}
