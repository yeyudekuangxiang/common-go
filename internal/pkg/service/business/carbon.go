package business

import (
	"errors"
	"fmt"
	ebusiness "mio/internal/pkg/model/entity/business"
	"mio/internal/pkg/util"
	"time"
)

type CarbonService struct {
}

// CarbonCreditEvCar 电测充电
func (srv CarbonService) CarbonCreditEvCar(userId int64, electricity int64) error {
	lockKey := fmt.Sprintf("CarbonCreditEvCar%d", userId)
	if !util.DefaultLock.Lock(lockKey, time.Second*10) {
		return errors.New("操作频率过快,请稍后再试")
	}
	defer util.DefaultLock.UnLock(lockKey)

	transactionId := util.UUID()
	_, _, ok, err := DefaultCarbonCreditsLimitService.CheckLimit(userId, ebusiness.CarbonTypeEvCar)
	if err != nil {
		return err
	}
	if ok {
		_, err = DefaultCarbonCreditsService.CarbonCreditEvCar(userId, electricity, transactionId)
	}

	//发放积分
	return nil
}

//CarbonCreditOnlineMeeting 在线会议
func (srv CarbonService) CarbonCreditOnlineMeeting(userId int64, duration time.Duration, start, end time.Time) error {
	lockKey := fmt.Sprintf("CarbonCreditOnlineMeeting%d", userId)
	if !util.DefaultLock.Lock(lockKey, time.Second*10) {
		return errors.New("操作频率过快,请稍后再试")
	}
	defer util.DefaultLock.UnLock(lockKey)

	transactionId := util.UUID()
	_, _, ok, err := DefaultCarbonCreditsLimitService.CheckLimit(userId, ebusiness.CarbonTypeEvCar)
	if err != nil {
		return err
	}
	if ok {
		_, err = DefaultCarbonCreditsService.CarbonCreditOnlineMeeting(userId, duration, start, end, transactionId)
	}

	//发放积分
	return nil
}

//CarbonCreditSaveWaterElectricity 节水节电
func (srv CarbonService) CarbonCreditSaveWaterElectricity(userId int64, water, electricity int64) error {
	lockKey := fmt.Sprintf("CarbonCreditSaveWaterElectricity%d", userId)
	if !util.DefaultLock.Lock(lockKey, time.Second*10) {
		return errors.New("操作频率过快,请稍后再试")
	}
	defer util.DefaultLock.UnLock(lockKey)

	transactionId := util.UUID()
	_, _, ok, err := DefaultCarbonCreditsLimitService.CheckLimit(userId, ebusiness.CarbonTypeEvCar)
	if err != nil {
		return err
	}
	if ok {
		_, err = DefaultCarbonCreditsService.CarbonCreditSaveWaterElectricity(userId, water, electricity, transactionId)
	}

	//发放积分
	return nil
}

//CarbonCreditSavePublicTransport 公交地铁
func (srv CarbonService) CarbonCreditSavePublicTransport(userId int64, bus int64, metro int64) error {
	lockKey := fmt.Sprintf("CarbonCreditSavePublicTransport%d", userId)
	if !util.DefaultLock.Lock(lockKey, time.Second*10) {
		return errors.New("操作频率过快,请稍后再试")
	}
	defer util.DefaultLock.UnLock(lockKey)

	transactionId := util.UUID()
	_, _, ok, err := DefaultCarbonCreditsLimitService.CheckLimit(userId, ebusiness.CarbonTypeEvCar)
	if err != nil {
		return err
	}
	if ok {
		_, err = DefaultCarbonCreditsService.CarbonCreditSavePublicTransport(userId, bus, metro, transactionId)
	}

	//发放积分
	return nil
}
