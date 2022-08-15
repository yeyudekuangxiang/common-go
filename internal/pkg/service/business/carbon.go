package business

import (
	"errors"
	"fmt"
	ebusiness "mio/internal/pkg/model/entity/business"
	"mio/internal/pkg/util"
	"time"
)

var DefaultCarbonService = CarbonService{}

type CarbonService struct {
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

	sendCarbonResult, err := DefaultCarbonCreditsService.SendCarbonCreditEvCar(SendCarbonCreditEvCarParam{
		UserId:        userId,
		Electricity:   electricity,
		TransactionId: transactionId,
	})
	if err != nil {
		return nil, err
	}

	//发送积分
	point, err := DefaultPointService.SendPointEvCar(SendPointEvCarParam{
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

	sendCarbonResult, err := DefaultCarbonCreditsService.SendCarbonCreditOnlineMeeting(SendCarbonCreditOnlineMeetingParam{
		UserId:           userId,
		OneCityDuration:  oneCityDuration,
		ManyCityDuration: manyCityDuration,
		TransactionId:    transactionId,
	})
	if err != nil {
		return nil, err
	}

	point, err := DefaultPointService.SendPointOnlineMeeting(SendPointOnlineMeetingParam{
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

	sendResult, err := DefaultCarbonCreditsService.SendCarbonCreditSaveWaterElectricity(SendCarbonCreditSaveWaterElectricityParam{
		UserId:        userId,
		Water:         water,
		Electricity:   electricity,
		TransactionId: transactionId,
	})
	if err != nil {
		return nil, err
	}

	point, err := DefaultPointService.SendPointSaveWaterElectricity(SendPointSaveWaterElectricityParam{
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

	sendResult, err := DefaultCarbonCreditsService.SendCarbonCreditSavePublicTransport(SendCarbonCreditSavePublicTransportParam{
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

	point, err := DefaultPointService.SendPointPublicTransport(SendPointPublicTransportParam{
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

	sendCarbonResult, err := DefaultCarbonCreditsService.SendCarbonCreditOEP(SendCarbonCreditOEPParam{
		UserId:        userId,
		Voucher:       voucher,
		TransactionId: transactionId,
	})
	if err != nil {
		return nil, err
	}

	//发送积分
	point, err := DefaultPointService.SendPointOEP(SendPointOEPParam{
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
func (srv CarbonService) CarbonCreditGreenBusinessTrip(userId int64, tripType string, from, to, voucher string) (*CarbonResult, error) {

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

	transactionId := util.UUID()

	sendCarbonResult, err := DefaultCarbonCreditsService.SendCarbonGreenBusinessTrip(SendCarbonGreenBusinessTripParam{
		TripType:      tripType,
		From:          from,
		To:            to,
		Voucher:       voucher,
		Distance:      0,
		UserId:        userId,
		TransactionId: transactionId,
	})
	if err != nil {
		return nil, err
	}

	//发送积分
	point, err := DefaultPointService.SendPointGreenBusinessTrip(SendPointGreenBusinessTripParam{
		TripType:      tripType,
		From:          from,
		To:            to,
		Voucher:       voucher,
		Distance:      0,
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
