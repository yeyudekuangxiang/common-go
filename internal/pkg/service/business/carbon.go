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
func (srv CarbonService) CarbonCreditPublicTransport(userId int64, bus float64, metro float64) (*CarbonResult, error) {
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
		TransactionId: transactionId,
	})
	if err != nil {
		return nil, err
	}

	point, err := DefaultPointService.SendPointPublicTransport(SendPointPublicTransportParam{
		UserId:        userId,
		Bus:           bus,
		Metro:         metro,
		BusCredit:     sendResult.BusCredits,
		MetroCredit:   sendResult.MetroCredits,
		TransactionId: transactionId,
	})
	if err != nil {
		return nil, err
	}

	return &CarbonResult{
		Credit: sendResult.BusCredits.Add(sendResult.MetroCredits),
		Point:  point,
	}, nil
}
