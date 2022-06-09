package business

import (
	"fmt"
	"mio/internal/pkg/model"
	ebusiness "mio/internal/pkg/model/entity/business"
	rbusiness "mio/internal/pkg/repository/business"
	"mio/internal/pkg/util"
	"time"
)

var DefaultPointService = PointService{repo: rbusiness.DefaultPointRepository}

type PointService struct {
	repo rbusiness.PointRepository
}

// SendPoint 发放碳积分 返回用户积分账户 本次实际发放的碳积分数量
func (srv PointService) SendPoint(param SendPointParam) (*ebusiness.Point, int, error) {
	lockKey := fmt.Sprintf("SendPoint%d", param.UserId)
	util.DefaultLock.LockWait(lockKey, time.Second*10)
	defer util.DefaultLock.UnLock(lockKey)

	//检测是否超过发放限制
	_, availableValue, err := DefaultPointLimitService.CheckLimitAndUpdate(param.UserId, param.AddPoint, param.Type)
	if err != nil {
		return nil, 0, err
	}

	//添加发放积分记录
	_, err = DefaultPointLogService.CreatePointLog(CreatePointLogParam{
		TransactionId: param.TransactionId,
		UserId:        param.UserId,
		Type:          param.Type,
		Value:         availableValue,
		Info:          param.Info,
		OrderId:       param.OrderId,
	})
	if err != nil {
		return nil, 0, err
	}

	//发放碳积分
	point, err := srv.createOrUpdatePoint(createOrUpdatePointParam{
		UserId:   param.UserId,
		AddPoint: availableValue,
	})
	if err != nil {
		return nil, 0, err
	}
	return point, availableValue, nil
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
		return &point, srv.repo.Save(&point)
	}

	point = ebusiness.Point{
		BUserId: param.UserId,
		Point:   int64(param.AddPoint),
	}
	return &point, srv.repo.Create(&point)
}

// PointEvCar 充电得碳积分
func (srv PointService) PointEvCar(userId int64, electricity float64, TransactionId string) (int, error) {

	//需要方法-查询用户信息
	userInfo := ebusiness.User{}
	sceneSetting, err := DefaultCompanyCarbonSceneService.FindCompanySceneSetting(userInfo.BCompanyId, ebusiness.CarbonTypeEvCar)
	if err != nil {
		return 0, err
	}

	addPoint := int(electricity/100) * sceneSetting.PointSetting

	_, point, err := srv.SendPoint(SendPointParam{
		UserId:        userId,
		AddPoint:      addPoint,
		Type:          ebusiness.PointTypeEvCar,
		TransactionId: TransactionId,
		Info: ebusiness.PointTypeInfo(ebusiness.CarbonTypeInfoEvCar{
			Electricity: electricity,
		}.JSON()),
	})
	return point, err
}

//PointOnlineMeeting 在线会议得碳积分
func (srv PointService) PointOnlineMeeting(userId int64, duration time.Duration, start, end time.Time, TransactionId string) (int, error) {
	//需要方法-查询用户信息
	userInfo := ebusiness.User{}
	_, err := DefaultCompanyCarbonSceneService.FindCompanySceneSetting(userInfo.BCompanyId, ebusiness.CarbonTypeOnlineMeeting)
	if err != nil {
		return 0, err
	}

	_, point, err := srv.SendPoint(SendPointParam{
		UserId:        userId,
		AddPoint:      0,
		Type:          ebusiness.PointTypeOnlineMeeting,
		TransactionId: TransactionId,
		Info: ebusiness.PointTypeInfo(ebusiness.CarbonTypeInfoOnlineMeeting{
			StartTime:       model.Time{Time: start},
			EndTime:         model.Time{Time: end},
			MeetingDuration: duration,
		}.JSON()),
	})
	return point, err
}

//PointSaveWaterElectricity 节水节电得积分
func (srv PointService) PointSaveWaterElectricity(userId int64, water, electricity int64, TransactionId string) (int, error) {
	//需要方法-查询用户信息
	userInfo := ebusiness.User{}
	_, err := DefaultCompanyCarbonSceneService.FindCompanySceneSetting(userInfo.BCompanyId, ebusiness.CarbonTypeSaveWaterElectricity)
	if err != nil {
		return 0, err
	}

	_, point, err := srv.SendPoint(SendPointParam{
		UserId:        userId,
		AddPoint:      0,
		Type:          ebusiness.PointTypeSaveWaterElectricity,
		TransactionId: TransactionId,
		Info: ebusiness.PointTypeInfo(ebusiness.CarbonTypeInfoSaveWaterElectricity{
			Water:       water,
			Electricity: electricity,
		}.JSON()),
	})
	return point, err
}

//PointPublicTransport 乘坐公交地铁得积分
func (srv PointService) PointPublicTransport(userId int64, bus int64, metro int64, TransactionId string) (int, error) {
	//需要方法-查询用户信息
	userInfo := ebusiness.User{}
	_, err := DefaultCompanyCarbonSceneService.FindCompanySceneSetting(userInfo.BCompanyId, ebusiness.CarbonTypePublicTransport)
	if err != nil {
		return 0, err
	}

	_, point, err := srv.SendPoint(SendPointParam{
		UserId:        userId,
		AddPoint:      0,
		Type:          ebusiness.PointTypePublicTransport,
		TransactionId: TransactionId,
		Info: ebusiness.PointTypeInfo(ebusiness.CarbonTypeInfoPublicTransport{
			Bus:   bus,
			Metro: metro,
		}.JSON()),
	})
	return point, err
}
