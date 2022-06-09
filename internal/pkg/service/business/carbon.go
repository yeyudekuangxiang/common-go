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

// CarbonCreditEvCar 电测充电
func (srv CarbonService) CarbonCreditEvCar(userId int64, electricity float64) error {
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

	_, _, ok, err = DefaultPointLimitService.CheckLimit(userId, ebusiness.PointTypeEvCar)
	if err != nil {
		return err
	}
	if ok {
		_, err = DefaultPointService.PointEvCar(userId, electricity, transactionId)
	}

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
	_, _, ok, err := DefaultCarbonCreditsLimitService.CheckLimit(userId, ebusiness.CarbonTypeOnlineMeeting)
	if err != nil {
		return err
	}
	if ok {
		_, err = DefaultCarbonCreditsService.CarbonCreditOnlineMeeting(userId, duration, start, end, transactionId)
	}

	_, _, ok, err = DefaultPointLimitService.CheckLimit(userId, ebusiness.PointTypeOnlineMeeting)
	if err != nil {
		return err
	}
	if ok {
		_, err = DefaultPointService.PointOnlineMeeting(userId, duration, start, end, transactionId)
	}
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
	_, _, ok, err := DefaultCarbonCreditsLimitService.CheckLimit(userId, ebusiness.CarbonTypeSaveWaterElectricity)
	if err != nil {
		return err
	}
	if ok {
		_, err = DefaultCarbonCreditsService.CarbonCreditSaveWaterElectricity(userId, water, electricity, transactionId)
	}

	_, _, ok, err = DefaultPointLimitService.CheckLimit(userId, ebusiness.PointTypeSaveWaterElectricity)
	if err != nil {
		return err
	}
	if ok {
		_, err = DefaultPointService.PointSaveWaterElectricity(userId, water, electricity, transactionId)
	}
	return nil
}

//CarbonCreditPublicTransport 公交地铁
func (srv CarbonService) CarbonCreditPublicTransport(userId int64, bus int64, metro int64) error {
	lockKey := fmt.Sprintf("CarbonCreditSavePublicTransport%d", userId)
	if !util.DefaultLock.Lock(lockKey, time.Second*10) {
		return errors.New("操作频率过快,请稍后再试")
	}
	defer util.DefaultLock.UnLock(lockKey)

	transactionId := util.UUID()
	_, _, ok, err := DefaultCarbonCreditsLimitService.CheckLimit(userId, ebusiness.CarbonTypePublicTransport)
	if err != nil {
		return err
	}
	if ok {
		_, err = DefaultCarbonCreditsService.CarbonCreditSavePublicTransport(userId, bus, metro, transactionId)
	}

	_, _, ok, err = DefaultPointLimitService.CheckLimit(userId, ebusiness.PointTypePublicTransport)
	if err != nil {
		return err
	}
	if ok {
		_, err = DefaultPointService.PointPublicTransport(userId, bus, metro, transactionId)
	}
	return nil
}
