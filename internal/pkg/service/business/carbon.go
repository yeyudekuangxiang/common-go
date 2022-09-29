package business

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	ebusiness "mio/internal/pkg/model/entity/business"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"strconv"
	"time"
)

type CarbonService struct {
	ctx *context.MioContext
}

func NewCarbonService(ctx *context.MioContext) *CarbonService {
	return &CarbonService{ctx: ctx}
}

// CarbonCreditEvCar 电车充电
func (srv CarbonService) CarbonCreditEvCar(userId int64, electricity float64) (*CarbonResult, error) {
	lockKey := fmt.Sprintf("CarbonCreditEvCar%d", userId)
	if !util.DefaultLock.Lock(lockKey, time.Second*10) {
		return nil, errors.New("操作频率过快,请稍后再试")
	}
	defer util.DefaultLock.UnLock(lockKey)

	//检测是否达到上限
	count, err := DefaultCarbonCreditsLimitService.CheckLimit(userId, ebusiness.CarbonTypeEvCar)
	if err != nil {
		return nil, err
	}
	if count <= 0 {
		return nil, errors.New("已经达到此场景当日最大限制")
	}

	transactionId := util.UUID()

	sendCarbonResult, err := NewCarbonCreditsService(srv.ctx).SendCarbonCreditEvCar(SendCarbonCreditEvCarParam{
		UserId:        userId,
		Electricity:   electricity,
		TransactionId: transactionId,
	})
	if err != nil {
		return nil, err
	}

	//发送积分
	point, err := NewPointService(srv.ctx).SendPointEvCar(SendPointEvCarParam{
		UserId:        userId,
		Electricity:   electricity,
		TransactionId: transactionId,
		CarbonCredits: sendCarbonResult.Credit,
	})

	if err != nil {
		return nil, err
	}

	return &CarbonResult{
		Credit: sendCarbonResult.Credit,
		Point:  point,
	}, nil
}

//CarbonCreditOnlineMeeting 在线会议
func (srv CarbonService) CarbonCreditOnlineMeeting(userId int64, oneCityDuration, manyCityDuration time.Duration) (*CarbonResult, error) {
	lockKey := fmt.Sprintf("CarbonCreditOnlineMeeting%d", userId)
	if !util.DefaultLock.Lock(lockKey, time.Second*10) {
		return nil, errors.New("操作频率过快,请稍后再试")
	}
	defer util.DefaultLock.UnLock(lockKey)

	transactionId := util.UUID()
	count, err := DefaultCarbonCreditsLimitService.CheckLimit(userId, ebusiness.CarbonTypeOnlineMeeting)
	if err != nil {
		return nil, err
	}
	if count <= 0 {
		return nil, errors.New("已经达到此场景当日最大限制")
	}

	sendCarbonResult, err := NewCarbonCreditsService(srv.ctx).SendCarbonCreditOnlineMeeting(SendCarbonCreditOnlineMeetingParam{
		UserId:           userId,
		OneCityDuration:  oneCityDuration,
		ManyCityDuration: manyCityDuration,
		TransactionId:    transactionId,
	})
	if err != nil {
		return nil, err
	}

	point, err := NewPointService(srv.ctx).SendPointOnlineMeeting(SendPointOnlineMeetingParam{
		UserId:           userId,
		OneCityDuration:  oneCityDuration,
		manyCityDuration: manyCityDuration,
		OneCityCredit:    sendCarbonResult.OneCityCredit,
		ManyCityCredit:   sendCarbonResult.ManyCityCredit,
		TransactionId:    transactionId,
	})

	if err != nil {
		return nil, err
	}
	return &CarbonResult{
		Credit: sendCarbonResult.OneCityCredit.Add(sendCarbonResult.ManyCityCredit),
		Point:  point,
	}, nil
}

//CarbonCreditSaveWaterElectricity 节水节电
func (srv CarbonService) CarbonCreditSaveWaterElectricity(userId int64, water, electricity int64) (*CarbonResult, error) {
	lockKey := fmt.Sprintf("CarbonCreditSaveWaterElectricity%d", userId)
	if !util.DefaultLock.Lock(lockKey, time.Second*10) {
		return nil, errors.New("操作频率过快,请稍后再试")
	}
	defer util.DefaultLock.UnLock(lockKey)

	transactionId := util.UUID()
	count, err := DefaultCarbonCreditsLimitService.CheckLimit(userId, ebusiness.CarbonTypeSaveWaterElectricity)
	if err != nil {
		return nil, err
	}
	if count <= 0 {
		return nil, errors.New("已经达到此场景当日最大限制")
	}

	sendResult, err := NewCarbonCreditsService(srv.ctx).SendCarbonCreditSaveWaterElectricity(SendCarbonCreditSaveWaterElectricityParam{
		UserId:        userId,
		Water:         water,
		Electricity:   electricity,
		TransactionId: transactionId,
	})
	if err != nil {
		return nil, err
	}

	point, err := NewPointService(srv.ctx).SendPointSaveWaterElectricity(SendPointSaveWaterElectricityParam{
		UserId:            userId,
		Water:             water,
		Electricity:       electricity,
		WaterCredit:       sendResult.WaterCredit,
		ElectricityCredit: sendResult.ElectricityCredit,
		TransactionId:     transactionId,
	})
	if err != nil {
		return nil, err
	}
	return &CarbonResult{
		Credit: sendResult.WaterCredit.Add(sendResult.ElectricityCredit),
		Point:  point,
	}, nil
}

//CarbonCreditPublicTransport 公交地铁
func (srv CarbonService) CarbonCreditPublicTransport(userId int64, bus, metro, step, bike float64) (*CarbonResult, error) {
	lockKey := fmt.Sprintf("CarbonCreditSavePublicTransport%d", userId)
	if !util.DefaultLock.Lock(lockKey, time.Second*10) {
		return nil, errors.New("操作频率过快,请稍后再试")
	}
	defer util.DefaultLock.UnLock(lockKey)

	transactionId := util.UUID()
	count, err := DefaultCarbonCreditsLimitService.CheckLimit(userId, ebusiness.CarbonTypePublicTransport)
	if err != nil {
		return nil, err
	}
	if count <= 0 {
		return nil, errors.New("已经达到此场景当日最大限制")
	}

	sendResult, err := NewCarbonCreditsService(srv.ctx).SendCarbonCreditSavePublicTransport(SendCarbonCreditSavePublicTransportParam{
		UserId:        userId,
		Bus:           bus,
		Metro:         metro,
		Step:          step,
		Bike:          bike,
		TransactionId: transactionId,
	})
	if err != nil {
		return nil, err
	}

	point, err := NewPointService(srv.ctx).SendPointPublicTransport(SendPointPublicTransportParam{
		UserId:        userId,
		Bus:           bus,
		Metro:         metro,
		Step:          step,
		Bike:          bike,
		BusCredit:     sendResult.BusCredits,
		MetroCredit:   sendResult.MetroCredits,
		StepCredit:    sendResult.StepCredits,
		BikeCredit:    sendResult.BikeCredits,
		TransactionId: transactionId,
	})
	if err != nil {
		return nil, err
	}

	return &CarbonResult{
		Credit: sendResult.BusCredits.Add(sendResult.MetroCredits).Add(sendResult.StepCredits).Add(sendResult.BikeCredits),
		Point:  point,
	}, nil
}

//CarbonCreditOEP 光盘行动
// userId 用户id
// voucher 凭证图片
func (srv CarbonService) CarbonCreditOEP(userId int64, voucher string) (*CarbonResult, error) {
	lockKey := fmt.Sprintf("CarbonCreditOEP%d", userId)
	if !util.DefaultLock.Lock(lockKey, time.Second*10) {
		return nil, errors.New("操作频率过快,请稍后再试")
	}
	defer util.DefaultLock.UnLock(lockKey)

	//检测是否达到上限
	count, err := DefaultCarbonCreditsLimitService.CheckLimit(userId, ebusiness.CarbonTypeOEP)
	if err != nil {
		return nil, err
	}
	if count <= 0 {
		return nil, errors.New("已经达到此场景当日最大限制")
	}

	transactionId := util.UUID()

	sendCarbonResult, err := NewCarbonCreditsService(srv.ctx).SendCarbonCreditOEP(SendCarbonCreditOEPParam{
		UserId:        userId,
		Voucher:       voucher,
		TransactionId: transactionId,
	})
	if err != nil {
		return nil, err
	}

	//发送积分
	point, err := NewPointService(srv.ctx).SendPointOEP(SendPointOEPParam{
		UserId:        userId,
		Voucher:       voucher,
		CarbonCredit:  sendCarbonResult.Credits,
		TransactionId: transactionId,
	})

	if err != nil {
		return nil, err
	}

	return &CarbonResult{
		Credit: sendCarbonResult.Credits,
		Point:  point,
	}, nil
}

// CarbonCreditGreenBusinessTrip 低碳出行
func (srv CarbonService) CarbonCreditGreenBusinessTrip(userId int64, tripType ebusiness.TripType, from, to, voucher string) (*CarbonResult, error) {

	lockKey := fmt.Sprintf("CarbonCreditGreenBusinessTrip%d", userId)
	if !util.DefaultLock.Lock(lockKey, time.Second*10) {
		return nil, errors.New("操作频率过快,请稍后再试")
	}
	defer util.DefaultLock.UnLock(lockKey)

	//检测是否达到上限
	count, err := DefaultCarbonCreditsLimitService.CheckLimit(userId, ebusiness.CarbonTypeGreenBusinessTrip)
	if err != nil {
		return nil, err
	}
	if count <= 0 {
		return nil, errors.New("已经达到此场景当日最大限制")
	}

	//计算城市距离
	distance, err := srv.CalcGreenBusinessTripCity(from, to)
	if err != nil {
		return nil, err
	}

	transactionId := util.UUID()
	sendCarbonResult, err := NewCarbonCreditsService(srv.ctx).SendCarbonGreenBusinessTrip(SendCarbonGreenBusinessTripParam{
		TripType:      tripType,
		From:          from,
		To:            to,
		Voucher:       voucher,
		Distance:      distance,
		UserId:        userId,
		TransactionId: transactionId,
	})
	if err != nil {
		return nil, err
	}

	//发送积分
	point, err := NewPointService(srv.ctx).SendPointGreenBusinessTrip(SendPointGreenBusinessTripParam{
		TripType:      tripType,
		From:          from,
		To:            to,
		Voucher:       voucher,
		Distance:      distance,
		UserId:        userId,
		TransactionId: transactionId,
		CarbonCredit:  sendCarbonResult.Credits,
	})

	if err != nil {
		return nil, err
	}

	return &CarbonResult{
		Credit: sendCarbonResult.Credits,
		Point:  point,
	}, nil
}

func (srv CarbonService) CalcGreenBusinessTripCity(from, to string) (decimal.Decimal, error) {
	areaSrv := NewAreaService(srv.ctx)

	fArea, exists, err := areaSrv.GetByName(ebusiness.AreaCity, from)
	if err != nil || !exists {
		return decimal.Decimal{}, errno.ErrCommon.WithMessage("出发城市名称错误")
	}

	tArea, exists, err := areaSrv.GetByName(ebusiness.AreaCity, to)
	if err != nil || !exists {
		return decimal.Decimal{}, errno.ErrCommon.WithMessage("到达城市名称错误")
	}

	app.Redis.ZIncrBy(srv.ctx, config.RedisKey.BusinessCarbonHotCity, 1, strconv.FormatInt(int64(fArea.CityID), 10))
	app.Redis.ZIncrBy(srv.ctx, config.RedisKey.BusinessCarbonHotCity, 1, strconv.FormatInt(int64(tArea.CityID), 10))
	lng1, err := strconv.ParseFloat(fArea.Longitude, 10)
	if err != nil {
		return decimal.Decimal{}, errno.ErrCommon.WithMessage("计算距离错误")
	}
	lat1, err := strconv.ParseFloat(fArea.Latitude, 10)
	if err != nil {
		return decimal.Decimal{}, errno.ErrCommon.WithMessage("计算距离错误")
	}

	lng2, err := strconv.ParseFloat(tArea.Longitude, 10)
	if err != nil {
		return decimal.Decimal{}, errno.ErrCommon.WithMessage("计算距离错误")
	}
	lat2, err := strconv.ParseFloat(tArea.Latitude, 10)
	if err != nil {
		return decimal.Decimal{}, errno.ErrCommon.WithMessage("计算距离错误")
	}
	distance := util.CalcLngLatDistance(lng1, lat1, lng2, lat2)
	return decimal.NewFromFloat(distance).Div(decimal.NewFromInt32(10000)).Round(2), nil
}

var defaultHotCity = map[int64]string{
	1560101035148619776: "北京市",
	1560101036503379968: "天津市",
	1560101119227637760: "上海市",
	1560101131009437696: "杭州市",
	1560101175020269568: "青岛市",
	1560101278741213184: "广州市",
	1560101283522719744: "深圳市",
	1559819477053370368: "香港",
}

func (srv CarbonService) GetCarbonHotCity() []ShortArea {
	hotCacheList, err := app.Redis.ZRevRangeWithScores(srv.ctx, config.RedisKey.BusinessCarbonHotCity, 50, -1).Result()
	if err != nil {
		app.Logger.Error(err)
	}

	areaSrv := NewAreaService(srv.ctx)
	hotCityIds := make([]int64, 0)

	for _, item := range hotCacheList {
		hotCityIds = append(hotCityIds, item.Member.(int64))
		if len(hotCityIds) == 8 {
			break
		}
	}

	for cityId := range defaultHotCity {
		if len(hotCityIds) == 8 {
			break
		}
		hotCityIds = append(hotCityIds, cityId)
	}

	hotList, err := areaSrv.List(AreaListDTO{
		CityIds: hotCityIds,
	})
	if err != nil {
		app.Logger.Error(err)
	}

	hotCityMap := make(map[int64]ebusiness.Area)
	for _, city := range hotList {
		hotCityMap[int64(city.CityID)] = city
	}

	areaList := make([]ShortArea, 0)

	for _, cityId := range hotCityIds {
		area := ShortArea{}
		err := util.MapTo(hotCityMap[cityId], &area)
		if err != nil {
			app.Logger.Errorf("map %+v to ShortArea{} err %+v", hotCityMap[cityId], err)
		}
		areaList = append(areaList, area)
	}

	return areaList
}