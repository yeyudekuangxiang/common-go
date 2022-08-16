package business

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api/business/businesstypes"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	ebusiness "mio/internal/pkg/model/entity/business"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/business"
	"mio/internal/pkg/util/apiutil"
	"mio/pkg/errno"
	"path/filepath"
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
	fh, err := ctx.FormFile("photo")
	if err != nil {
		app.Logger.Error(err)
		return nil, errno.ErrCommon.WithMessage("上传图片失败,请稍后再试")
	}
	f, err := fh.Open()
	if err != nil {
		app.Logger.Error(err)
		return nil, errno.ErrCommon.WithMessage("上传图片失败,请稍后再试")
	}
	defer f.Close()

	user := apiutil.GetAuthBusinessUser(ctx)
	path := fmt.Sprintf("business/carbon/oep/%d%d%s", user.ID, time.Now().UnixMilli(), filepath.Ext(fh.Filename))
	path, err = service.DefaultOssService.PutObject(path, f)
	if err != nil {
		app.Logger.Error(err)
		return nil, errno.ErrCommon.WithMessage("上传图片失败,请稍后再试")
	}

	carbonSrv := business.NewCarbonService(context.NewMioContext(context.WithContext(ctx)))
	result, err := carbonSrv.CarbonCreditOEP(user.ID, path)
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

	fh, err := ctx.FormFile("photo")
	if err != nil {
		app.Logger.Error(err)
		return nil, errno.ErrCommon.WithMessage("上传图片失败,请稍后再试")
	}
	f, err := fh.Open()
	if err != nil {
		app.Logger.Error(err)
		return nil, errno.ErrCommon.WithMessage("上传图片失败,请稍后再试")
	}
	defer f.Close()

	user := apiutil.GetAuthBusinessUser(ctx)
	path := fmt.Sprintf("business/carbon/gbt/%d%d%s", user.ID, time.Now().UnixMilli(), filepath.Ext(fh.Filename))
	path, err = service.DefaultOssService.PutObject(path, f)
	if err != nil {
		app.Logger.Error(err)
		return nil, errno.ErrCommon.WithMessage("上传图片失败,请稍后再试")
	}

	carbonSrv := business.NewCarbonService(context.NewMioContext(context.WithContext(ctx)))
	result, err := carbonSrv.CarbonCreditGreenBusinessTrip(user.ID, ebusiness.TripType(form.Type), form.From, form.To, path)
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
