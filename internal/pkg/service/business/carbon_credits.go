package business

import (
	"fmt"
	"github.com/shopspring/decimal"
	"mio/internal/pkg/model"
	ebusiness "mio/internal/pkg/model/entity/business"
	rbusiness "mio/internal/pkg/repository/business"
	"mio/internal/pkg/util"
	"time"
)

var DefaultCarbonCreditsService = CarbonCreditsService{repo: rbusiness.DefaultCarbonCreditsRepository}

type CarbonCreditsService struct {
	repo rbusiness.CarbonCreditsRepository
}

// SendCarbonCredit 发放碳积分 返回用户积分账户 本次实际发放的碳积分数量
func (srv CarbonCreditsService) SendCarbonCredit(param SendCarbonCreditParam) (*ebusiness.CarbonCredits, decimal.Decimal, error) {
	lockKey := fmt.Sprintf("SendCarbonCredit%d", param.UserId)
	util.DefaultLock.LockWait(lockKey, time.Second*10)
	defer util.DefaultLock.UnLock(lockKey)

	//检测是否超过发放限制
	_, availableValue, err := DefaultCarbonCreditsLimitService.CheckLimitAndUpdate(param.UserId, param.AddCredit, param.Type)
	if err != nil {
		return nil, decimal.Decimal{}, err
	}

	//添加发放碳积分记录
	_, err = DefaultCarbonCreditsLogService.CreateCarbonCreditLog(CreateCarbonCreditLogParam{
		TransactionId: param.TransactionId,
		UserId:        param.UserId,
		Type:          param.Type,
		Value:         availableValue,
		Info:          param.Info,
	})
	if err != nil {
		return nil, decimal.Decimal{}, err
	}

	//发放碳积分
	carbonCredit, err := srv.createOrUpdateCarbonCredit(createOrUpdateCarbonCreditParam{
		UserId:    param.UserId,
		AddCredit: availableValue,
	})
	if err != nil {
		return nil, decimal.Decimal{}, err
	}
	return carbonCredit, availableValue, nil
}

//创建或者更新用户碳积分账户
func (srv CarbonCreditsService) createOrUpdateCarbonCredit(param createOrUpdateCarbonCreditParam) (*ebusiness.CarbonCredits, error) {
	lockKey := fmt.Sprintf("CreateOrUpdateCarbonCredit%d", param.UserId)
	util.DefaultLock.LockWait(lockKey, time.Second*10)
	defer util.DefaultLock.UnLock(lockKey)

	credit := srv.repo.FindCredits(param.UserId)
	if credit.ID != 0 {
		credit.Credits = credit.Credits.Add(param.AddCredit)
		return &credit, srv.repo.Save(&credit)
	}

	credit = ebusiness.CarbonCredits{
		BUserId: param.UserId,
		Credits: param.AddCredit,
	}
	return &credit, srv.repo.Create(&credit)
}

// CarbonCreditEvCar 充电得碳积分
func (srv CarbonCreditsService) CarbonCreditEvCar(userId int64, electricity int64, TransactionId string) (decimal.Decimal, error) {
	credits := ebusiness.CarbonTypeEvCar.CalcEvCar(electricity)
	_, credits, err := srv.SendCarbonCredit(SendCarbonCreditParam{
		UserId:        userId,
		AddCredit:     credits,
		Type:          ebusiness.CarbonTypeEvCar,
		TransactionId: TransactionId,
		Info: ebusiness.CarbonTypeInfoEvCar{
			Electricity: electricity,
		}.JSON(),
	})
	return credits, err
}

//CarbonCreditOnlineMeeting 在线会议得碳积分
func (srv CarbonCreditsService) CarbonCreditOnlineMeeting(userId int64, duration time.Duration, start, end time.Time, TransactionId string) (decimal.Decimal, error) {
	credits := ebusiness.CarbonTypeOnlineMeeting.CalcOnlineMeeting(duration)

	_, credits, err := srv.SendCarbonCredit(SendCarbonCreditParam{
		UserId:        userId,
		AddCredit:     credits,
		Type:          ebusiness.CarbonTypeEvCar,
		TransactionId: TransactionId,
		Info: ebusiness.CarbonTypeInfoOnlineMeeting{
			MeetingDuration: duration,
			StartTime:       model.Time{Time: start},
			EndTime:         model.Time{Time: end},
		}.JSON(),
	})
	return credits, err
}

//CarbonCreditSaveWaterElectricity 节水节电得积分
func (srv CarbonCreditsService) CarbonCreditSaveWaterElectricity(userId int64, water, electricity int64, TransactionId string) (decimal.Decimal, error) {
	credits := ebusiness.CarbonTypeSaveWaterElectricity.CalcSaveWaterElectricity(water, electricity)

	_, credits, err := srv.SendCarbonCredit(SendCarbonCreditParam{
		UserId:        userId,
		AddCredit:     credits,
		Type:          ebusiness.CarbonTypeEvCar,
		TransactionId: TransactionId,
		Info: ebusiness.CarbonTypeInfoSaveWaterElectricity{
			Water:       water,
			Electricity: electricity,
		}.JSON(),
	})
	return credits, err
}

//CarbonCreditSavePublicTransport 乘坐公交地铁得积分
func (srv CarbonCreditsService) CarbonCreditSavePublicTransport(userId int64, bus int64, metro int64, TransactionId string) (decimal.Decimal, error) {
	credits := ebusiness.CarbonTypePublicTransport.CalcPublicTransport(bus, metro)

	_, credits, err := srv.SendCarbonCredit(SendCarbonCreditParam{
		UserId:        userId,
		AddCredit:     credits,
		Type:          ebusiness.CarbonTypeEvCar,
		TransactionId: TransactionId,
		Info: ebusiness.CarbonTypeInfoPublicTransport{
			Bus:   bus,
			Metro: metro,
		}.JSON(),
	})
	return credits, err
}
