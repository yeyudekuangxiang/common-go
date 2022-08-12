package business

import (
	"fmt"
	"github.com/shopspring/decimal"
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
	_, err := DefaultCarbonCreditsLimitService.CheckLimitAndUpdate(param.UserId, param.AddCredit, param.Type)
	if err != nil {
		return nil, decimal.Decimal{}, err
	}

	//添加发放碳积分记录
	_, err = DefaultCarbonCreditsLogService.CreateCarbonCreditLog(CreateCarbonCreditLogParam{
		TransactionId: param.TransactionId,
		UserId:        param.UserId,
		Type:          param.Type,
		Value:         param.AddCredit,
		Info:          param.Info,
	})
	if err != nil {
		return nil, decimal.Decimal{}, err
	}

	//发放碳积分
	carbonCredit, err := srv.createOrUpdateCarbonCredit(createOrUpdateCarbonCreditParam{
		UserId:    param.UserId,
		AddCredit: param.AddCredit,
	})
	if err != nil {
		return nil, decimal.Decimal{}, err
	}
	return carbonCredit, param.AddCredit, nil
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

// SendCarbonCreditEvCar 充电得碳积分
func (srv CarbonCreditsService) SendCarbonCreditEvCar(param SendCarbonCreditEvCarParam) (*CarbonCreditEvCarResult, error) {
	//计算碳积分
	credits := DefaultCarbonCreditCalculatorService.CalcEvCar(param.Electricity)

	//发放碳积分
	_, credits, err := DefaultCarbonCreditsService.SendCarbonCredit(SendCarbonCreditParam{
		UserId:        param.UserId,
		AddCredit:     credits,
		Type:          ebusiness.CarbonTypeEvCar,
		TransactionId: param.TransactionId,
		Info: ebusiness.CarbonTypeInfoEvCar{
			Electricity: param.Electricity,
		}.CarbonTypeInfo(),
	})
	if err != nil {
		return nil, err
	}
	return &CarbonCreditEvCarResult{Credit: credits}, nil
}

//SendCarbonCreditOnlineMeeting 在线会议得碳积分
func (srv CarbonCreditsService) SendCarbonCreditOnlineMeeting(param SendCarbonCreditOnlineMeetingParam) (*SendCarbonCreditOnlineMeetingResult, error) {
	oneCityCredit := DefaultCarbonCreditCalculatorService.CalcOnlineMeetingOneCity(param.OneCityDuration)
	manyCityCredit := DefaultCarbonCreditCalculatorService.CalcOnlineMeetingManyCity(param.ManyCityDuration)

	_, _, err := srv.SendCarbonCredit(SendCarbonCreditParam{
		UserId:        param.UserId,
		AddCredit:     oneCityCredit.Add(manyCityCredit),
		Type:          ebusiness.CarbonTypeOnlineMeeting,
		TransactionId: param.TransactionId,
		Info: ebusiness.CarbonTypeInfoOnlineMeeting{
			OneCityDuration:  param.OneCityDuration,
			ManyCityDuration: param.ManyCityDuration,
		}.CarbonTypeInfo(),
	})
	if err != nil {
		return nil, err
	}
	return &SendCarbonCreditOnlineMeetingResult{
		OneCityCredit:  oneCityCredit,
		ManyCityCredit: manyCityCredit,
	}, nil
}

//SendCarbonCreditSaveWaterElectricity 节水节电得积分
func (srv CarbonCreditsService) SendCarbonCreditSaveWaterElectricity(param SendCarbonCreditSaveWaterElectricityParam) (*SendCarbonCreditSaveWaterElectricityResult, error) {
	waterCredit := DefaultCarbonCreditCalculatorService.CalcSaveWater(param.Water)
	electricityCredit := DefaultCarbonCreditCalculatorService.CalcSaveElectricity(param.Electricity)
	_, _, err := srv.SendCarbonCredit(SendCarbonCreditParam{
		UserId:        param.UserId,
		AddCredit:     waterCredit.Add(electricityCredit),
		Type:          ebusiness.CarbonTypeSaveWaterElectricity,
		TransactionId: param.TransactionId,
		Info: ebusiness.CarbonTypeInfoSaveWaterElectricity{
			Water:       param.Water,
			Electricity: param.Electricity,
		}.CarbonTypeInfo(),
	})
	if err != nil {
		return nil, err
	}
	return &SendCarbonCreditSaveWaterElectricityResult{
		WaterCredit:       waterCredit,
		ElectricityCredit: electricityCredit,
	}, nil
}

//SendCarbonCreditSavePublicTransport 乘坐公交地铁得积分
func (srv CarbonCreditsService) SendCarbonCreditSavePublicTransport(param SendCarbonCreditSavePublicTransportParam) (*SendCarbonCreditSavePublicTransportResult, error) {
	busCredits := DefaultCarbonCreditCalculatorService.CalcBus(param.Bus)
	metroCredits := DefaultCarbonCreditCalculatorService.CalcMetro(param.Metro)

	_, _, err := srv.SendCarbonCredit(SendCarbonCreditParam{
		UserId:        param.UserId,
		AddCredit:     busCredits.Add(metroCredits),
		Type:          ebusiness.CarbonTypePublicTransport,
		TransactionId: param.TransactionId,
		Info: ebusiness.CarbonTypeInfoPublicTransport{
			Bus:   param.Bus,
			Metro: param.Metro,
		}.CarbonTypeInfo(),
	})
	if err != nil {
		return nil, err
	}
	return &SendCarbonCreditSavePublicTransportResult{
		BusCredits:   busCredits,
		MetroCredits: metroCredits,
	}, err
}

// SendCarbonCreditOEP 光盘行动得积分
func (srv CarbonCreditsService) SendCarbonCreditOEP(param SendCarbonCreditOEPParam) (*SendCarbonCreditOEPResult, error) {
	oepCredits := DefaultCarbonCreditCalculatorService.CalcOEP()
	_, _, err := srv.SendCarbonCredit(SendCarbonCreditParam{
		UserId:        param.UserId,
		AddCredit:     oepCredits,
		Type:          ebusiness.CarbonTypeOEP,
		TransactionId: param.TransactionId,
		Info: ebusiness.CarbonTypeInfoOEP{
			Voucher: param.Voucher,
		}.CarbonTypeInfo(),
	})
	if err != nil {
		return nil, err
	}
	return &SendCarbonCreditOEPResult{
		Credits: oepCredits,
	}, err
}
