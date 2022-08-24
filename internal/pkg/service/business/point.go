package business

import (
	"errors"
	"fmt"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	ebusiness "mio/internal/pkg/model/entity/business"
	rbusiness "mio/internal/pkg/repository/business"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"time"
)

type PointService struct {
	ctx  *context.MioContext
	repo *rbusiness.PointRepository
}

func NewPointService(ctx *context.MioContext) *PointService {
	return &PointService{ctx: ctx, repo: rbusiness.NewPointRepository(ctx)}
}

// SendPoint 发放碳积分 返回用户积分账户 本次实际发放的碳积分数量
func (srv PointService) SendPoint(param SendPointParam) (*ebusiness.Point, error) {
	lockKey := fmt.Sprintf("SendPoint%d", param.UserId)
	util.DefaultLock.LockWait(lockKey, time.Second*10)
	defer util.DefaultLock.UnLock(lockKey)

	//添加发放积分记录
	_, err := DefaultPointLogService.CreatePointLog(CreatePointLogParam{
		TransactionId: param.TransactionId,
		UserId:        param.UserId,
		Type:          param.Type,
		Value:         param.AddPoint,
		Info:          param.Info,
		OrderId:       param.OrderId,
	})
	if err != nil {
		return nil, err
	}

	//发放碳积分
	point, err := srv.createOrUpdatePoint(createOrUpdatePointParam{
		UserId:   param.UserId,
		AddPoint: param.AddPoint,
	})
	if err != nil {
		return nil, err
	}
	return point, nil
}

//创建或者更新用户碳积分账户
func (srv PointService) createOrUpdatePoint(param createOrUpdatePointParam) (*ebusiness.Point, error) {
	lockKey := fmt.Sprintf("CreateOrUpdatePoint%d", param.UserId)
	util.DefaultLock.LockWait(lockKey, time.Second*10)
	defer util.DefaultLock.UnLock(lockKey)

	point := srv.repo.FindPoint(rbusiness.FindPointBy{
		UserId: param.UserId,
	})

	if point.ID != 0 {
		point.Point += int64(param.AddPoint)
		srv.trackPointChange(param.UserId, param.AddPoint)
		return &point, srv.repo.Save(&point)
	}

	point = ebusiness.Point{
		BUserId: param.UserId,
		Point:   int64(param.AddPoint),
	}
	err := srv.repo.Create(&point)
	if err != nil {
		return nil, err
	}
	srv.trackPointChange(param.UserId, param.AddPoint)
	return &point, nil
}

// SendPointEvCar 充电得碳积分
func (srv PointService) SendPointEvCar(param SendPointEvCarParam) (int, error) {

	userInfo, err := DefaultUserService.GetBusinessUserById(param.UserId)
	if err != nil {
		return 0, err
	}
	if userInfo.ID == 0 {
		return 0, errno.ErrUserNotFound
	}

	sceneSetting, err := DefaultCompanyCarbonSceneService.FindCompanySceneSetting(userInfo.BCompanyId, ebusiness.CarbonTypeEvCar)
	if err != nil {
		return 0, err
	}

	//获取碳积分和积分汇率
	pointRate, err := DefaultPointRateSettingService.ParsePointExchangeRate(sceneSetting.PointRateSetting)
	if err != nil {
		app.Logger.Error("转换碳积分汇率异常", sceneSetting.PointRateSetting, err)
		return 0, errors.New("系统异常,请稍后再试")
	}

	addPoint := pointRate.Calc(param.CarbonCredits)

	_, err = srv.SendPoint(SendPointParam{
		UserId:        param.UserId,
		AddPoint:      addPoint,
		Type:          ebusiness.PointTypeEvCar,
		TransactionId: param.TransactionId,
		Info: ebusiness.CarbonTypeInfoEvCar{
			Electricity: param.Electricity,
		}.PointTypeInfo(),
	})
	return addPoint, err
}

//SendPointOnlineMeeting 在线会议得碳积分
func (srv PointService) SendPointOnlineMeeting(param SendPointOnlineMeetingParam) (int, error) {

	userInfo, err := DefaultUserService.GetBusinessUserById(param.UserId)
	if err != nil {
		return 0, err
	}
	if userInfo.ID == 0 {
		return 0, errno.ErrUserNotFound
	}

	sceneSetting, err := DefaultCompanyCarbonSceneService.FindCompanySceneSetting(userInfo.BCompanyId, ebusiness.CarbonTypeOnlineMeeting)
	if err != nil {
		return 0, err
	}

	//获取碳积分和积分汇率
	pointRate, err := DefaultPointRateSettingService.ParseOnlineMeetingRate(sceneSetting.PointRateSetting)
	if err != nil {
		app.Logger.Error("转换碳积分汇率异常", sceneSetting.PointRateSetting, err)
		return 0, errors.New("系统异常,请稍后再试")
	}

	addPoint := pointRate.OneCity.Calc(param.OneCityCredit)
	addPoint += pointRate.ManyCity.Calc(param.ManyCityCredit)

	_, err = srv.SendPoint(SendPointParam{
		UserId:        param.UserId,
		AddPoint:      addPoint,
		Type:          ebusiness.PointTypeOnlineMeeting,
		TransactionId: param.TransactionId,
		Info: ebusiness.CarbonTypeInfoOnlineMeeting{
			OneCityDuration:  param.OneCityDuration,
			ManyCityDuration: param.manyCityDuration,
		}.PointTypeInfo(),
	})
	return addPoint, err
}

//SendPointSaveWaterElectricity 节水节电得积分
func (srv PointService) SendPointSaveWaterElectricity(param SendPointSaveWaterElectricityParam) (int, error) {
	userInfo, err := DefaultUserService.GetBusinessUserById(param.UserId)
	if err != nil {
		return 0, err
	}
	if userInfo.ID == 0 {
		return 0, errno.ErrUserNotFound
	}

	sceneSetting, err := DefaultCompanyCarbonSceneService.FindCompanySceneSetting(userInfo.BCompanyId, ebusiness.CarbonTypeSaveWaterElectricity)
	if err != nil {
		return 0, err
	}

	//获取碳积分和积分汇率
	pointRate, err := DefaultPointRateSettingService.ParseSaveWaterElectricityRate(sceneSetting.PointRateSetting)
	if err != nil {
		app.Logger.Error("转换碳积分汇率异常", sceneSetting.PointRateSetting, err)
		return 0, errors.New("系统异常,请稍后再试")
	}
	addPoint := pointRate.Water.Calc(param.WaterCredit)
	addPoint += pointRate.Electricity.Calc(param.ElectricityCredit)

	_, err = srv.SendPoint(SendPointParam{
		UserId:        param.UserId,
		AddPoint:      addPoint,
		Type:          ebusiness.PointTypeSaveWaterElectricity,
		TransactionId: param.TransactionId,
		Info: ebusiness.CarbonTypeInfoSaveWaterElectricity{
			Water:       param.Water,
			Electricity: param.Electricity,
		}.PointTypeInfo(),
	})
	return addPoint, err
}

//SendPointPublicTransport 乘坐公交地铁得积分
func (srv PointService) SendPointPublicTransport(param SendPointPublicTransportParam) (int, error) {
	userInfo, err := DefaultUserService.GetBusinessUserById(param.UserId)
	if err != nil {
		return 0, err
	}
	if userInfo.ID == 0 {
		return 0, errno.ErrUserNotFound
	}

	sceneSetting, err := DefaultCompanyCarbonSceneService.FindCompanySceneSetting(userInfo.BCompanyId, ebusiness.CarbonTypePublicTransport)
	if err != nil {
		return 0, err
	}

	//获取碳积分和积分汇率
	pointRate, err := DefaultPointRateSettingService.ParsePublicTransportRate(sceneSetting.PointRateSetting)
	if err != nil {
		app.Logger.Error("转换碳积分汇率异常", sceneSetting.PointRateSetting, err)
		return 0, errors.New("系统异常,请稍后再试")
	}

	addPoint := pointRate.Bus.Calc(param.BusCredit)
	addPoint += pointRate.Metro.Calc(param.MetroCredit)
	addPoint += pointRate.Step.Calc(param.StepCredit)
	addPoint += pointRate.Bike.Calc(param.BikeCredit)

	_, err = srv.SendPoint(SendPointParam{
		UserId:        param.UserId,
		AddPoint:      addPoint,
		Type:          ebusiness.PointTypePublicTransport,
		TransactionId: param.TransactionId,
		Info: ebusiness.CarbonTypeInfoPublicTransport{
			Bus:   param.Bus,
			Metro: param.Metro,
			Step:  param.Step,
			Bike:  param.Bike,
		}.PointTypeInfo(),
	})
	return addPoint, err
}

func (srv PointService) SendPointOEP(param SendPointOEPParam) (int, error) {
	userInfo, err := DefaultUserService.GetBusinessUserById(param.UserId)
	if err != nil {
		return 0, err
	}
	if userInfo.ID == 0 {
		return 0, errno.ErrUserNotFound
	}

	sceneSetting, err := DefaultCompanyCarbonSceneService.FindCompanySceneSetting(userInfo.BCompanyId, ebusiness.CarbonTypeOEP)
	if err != nil {
		return 0, err
	}

	//获取碳积分和积分汇率
	pointRate, err := DefaultPointRateSettingService.ParsePointOEPRate(sceneSetting.PointRateSetting)
	if err != nil {
		app.Logger.Error("转换碳积分汇率异常", sceneSetting.PointRateSetting, err)
		return 0, errors.New("系统异常,请稍后再试")
	}

	addPoint := pointRate.Calc(param.CarbonCredit)

	_, err = srv.SendPoint(SendPointParam{
		UserId:        param.UserId,
		AddPoint:      addPoint,
		Type:          ebusiness.PointTypeOEP,
		TransactionId: param.TransactionId,
		Info: ebusiness.CarbonTypeInfoOEP{
			Voucher: param.Voucher,
		}.PointTypeInfo(),
	})
	return addPoint, err
}
func (srv PointService) SendPointGreenBusinessTrip(param SendPointGreenBusinessTripParam) (int, error) {
	userInfo, err := DefaultUserService.GetBusinessUserById(param.UserId)
	if err != nil {
		return 0, err
	}
	if userInfo.ID == 0 {
		return 0, errno.ErrUserNotFound
	}

	sceneSetting, err := DefaultCompanyCarbonSceneService.FindCompanySceneSetting(userInfo.BCompanyId, ebusiness.CarbonTypeGreenBusinessTrip)
	if err != nil {
		return 0, err
	}

	//获取碳积分和积分汇率
	pointRate, err := DefaultPointRateSettingService.ParseGreenBusinessTripExchangeRate(sceneSetting.PointRateSetting)
	if err != nil {
		app.Logger.Error("转换碳积分汇率异常", sceneSetting.PointRateSetting, err)
		return 0, errors.New("系统异常,请稍后再试")
	}
	var addPoint int
	switch param.TripType {
	case ebusiness.TripTypeTrain:
		addPoint = pointRate.Train.Calc(param.CarbonCredit)
	case ebusiness.TripTypeHighSpeed:
		addPoint = pointRate.HighSpeed.Calc(param.CarbonCredit)
	case ebusiness.TripTypeAirPlane:
		addPoint = pointRate.Airplane.Calc(param.CarbonCredit)
	}

	_, err = srv.SendPoint(SendPointParam{
		UserId:        param.UserId,
		AddPoint:      addPoint,
		Type:          ebusiness.PointTypeGreenBusinessTrip,
		TransactionId: param.TransactionId,
		Info: ebusiness.CarbonTypeInfoGreenBusinessTrip{
			TripType: param.TripType,
			Distance: param.Distance,
			From:     param.From,
			To:       param.To,
			Voucher:  param.Voucher,
		}.PointTypeInfo(),
	})
	return addPoint, err
}

func (srv PointService) trackPointChange(userId int64, value int) {
	go func() {
		userInfo, err := DefaultUserService.GetBusinessUserById(userId)
		if err != nil {
			return
		}
		department, err := DefaultDepartmentService.GetBusinessDepartmentById(userInfo.BDepartmentId)
		if err != nil {
			return
		}
		company := DefaultCompanyService.GetCompanyById(userInfo.BCompanyId)

		service.DefaultZhuGeService().TrackBusinessPoints(srv_types.TrackBusinessPoints{
			Uid:        userInfo.Uid,
			Value:      value,
			ChangeType: util.Ternary(value > 0, "inc", "dec").String(),
			Nickname:   userInfo.Nickname,
			Username:   userInfo.Realname,
			Department: department.Title,
			Company:    company.Name,
			ChangeTime: time.Now(),
		})
	}()
}
