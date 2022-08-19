package business

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api/business/businesstypes"
	"mio/internal/pkg/core/context"
	ebusiness "mio/internal/pkg/model/entity/business"
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
	carbonSrv := business.NewCarbonService(context.NewMioContext(context.WithContext(ctx)))
	result, err := carbonSrv.CarbonCreditEvCar(user.ID, form.Electricity)
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
	carbonSrv := business.NewCarbonService(context.NewMioContext(context.WithContext(ctx)))
	result, err := carbonSrv.CarbonCreditOnlineMeeting(user.ID, oneCityDuration, manyCityDuration)
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

	carbonSrv := business.NewCarbonService(context.NewMioContext(context.WithContext(ctx)))
	result, err := carbonSrv.CarbonCreditSaveWaterElectricity(user.ID, form.Water, form.Electricity)
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
	carbonSrv := business.NewCarbonService(context.NewMioContext(context.WithContext(ctx)))
	result, err := carbonSrv.CarbonCreditPublicTransport(user.ID, form.Bus, form.Metro, form.Walk, form.Bike)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"carbonCredit": result.Credit,
		"point":        result.Point,
	}, err
}
func (CarbonController) CollectOEP(ctx *gin.Context) (gin.H, error) {
	form := CarbonCollectOEPForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthBusinessUser(ctx)

	carbonSrv := business.NewCarbonService(context.NewMioContext(context.WithContext(ctx)))
	result, err := carbonSrv.CarbonCreditOEP(user.ID, form.Photo)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"carbonCredit": result.Credit,
		"point":        result.Point,
	}, err
}
func (CarbonController) CollectGreenBusinessTrip(ctx *gin.Context) (gin.H, error) {
	form := CarbonGreenBusinessTripForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthBusinessUser(ctx)
	carbonSrv := business.NewCarbonService(context.NewMioContext(context.WithContext(ctx)))
	result, err := carbonSrv.CarbonCreditGreenBusinessTrip(user.ID, ebusiness.TripType(form.Type), form.From, form.To, form.Photo)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"carbonCredit": result.Credit,
		"point":        result.Point,
	}, err
}
func (CarbonController) CityProvinceList(ctx *gin.Context) (gin.H, error) {
	form := businesstypes.CityProvinceForm{}

	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	areaSrv := business.NewAreaService(context.NewMioContext(context.WithContext(ctx)))
	list, err := areaSrv.GroupCityProvinceList(business.CityProvinceListDTO{
		Search: form.Search,
	})
	if err != nil {
		return nil, err
	}

	carbonSrv := business.NewCarbonService(context.NewMioContext(context.WithContext(ctx)))

	return gin.H{
		"hotCities": carbonSrv.GetCarbonHotCity(),
		"group":     list,
	}, nil
}
