package business

import (
	"errors"
	"fmt"
	ebusiness "mio/internal/pkg/model/entity/business"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/timeutils"
	"mio/pkg/errno"
	"time"
)

var DefaultPointLimitService = PointLimitService{}

type PointLimitService struct {
}

//检测是否超出限制 返回还可领取次数 还可领取积分值 是否超出
func (srv PointLimitService) checkLimit(userId int64, t ebusiness.PointType) (int, int, bool, error) {

	carbonType := t.CarbonType()
	if carbonType == "" {
		return 1, 9999, true, nil
	}

	userInfo, err := DefaultUserService.GetBusinessUserById(userId)
	if err != nil {
		return 0, 0, false, err
	}
	if userInfo.ID == 0 {
		return 0, 0, false, errno.ErrUserNotFound
	}

	carbonScene, err := DefaultCarbonSceneService.FindScene(carbonType)
	if err != nil {
		return 0, 0, false, err
	}
	if carbonScene.ID == 0 {
		return 0, 0, false, errors.New("未查询到此低碳场景")
	}
	companyCarbonScene, err := DefaultCompanyCarbonSceneService.FindCompanyScene(FindCompanyCarbonSceneParam{
		CompanyId: userInfo.BCompanyId,
	})
	if err != nil {
		return 0, 0, false, err
	}
	if companyCarbonScene.ID == 0 {
		return 0, 0, false, errors.New("未查询到此低碳场景")
	}

	limitLog, err := DefaultPointLimitLogService.FindLimitLog(FindPointLimitLogParam{
		TimePoint: timeutils.Now().StartOfDay().Time(),
		Type:      t,
		UserId:    userId,
	})
	if err != nil {
		return 0, 0, false, err
	}

	if limitLog.ID == 0 {
		return companyCarbonScene.MaxCount, carbonScene.MaxPoint, true, nil
	}

	count := companyCarbonScene.MaxCount - limitLog.CurrentCount
	if count <= 0 {
		return 0, 0, false, nil
	}

	credits := carbonScene.MaxPoint - limitLog.CurrentValue
	if credits <= 0 {
		return 0, 0, false, nil
	}

	return count, credits, true, nil
}

// CheckLimit 加读写锁检测是否超出限制 返回还可领取次数 还可领取积分值 是否超出
func (srv PointLimitService) CheckLimit(userId int64, t ebusiness.PointType) (int, int, bool, error) {
	lockKey := fmt.Sprintf("CheckPointLimit%d%s", userId, t)
	util.DefaultLock.LockWait(lockKey, time.Second*5)
	defer util.DefaultLock.UnLock(lockKey)
	return srv.checkLimit(userId, t)
}

//CheckLimitAndUpdate 加读写锁检测是否超出限制并且修改值 返回记录 以及实际新增的积分数量
func (srv PointLimitService) CheckLimitAndUpdate(userId int64, value int, t ebusiness.PointType) (*ebusiness.PointLimitLog, int, error) {
	lockKey := fmt.Sprintf("CheckPointLimit%d%s", userId, t)
	util.DefaultLock.LockWait(lockKey, time.Second*5)
	defer util.DefaultLock.UnLock(lockKey)

	//不是低碳场景不需要检查
	if t.CarbonType() == "" {
		return &ebusiness.PointLimitLog{}, value, nil
	}

	_, leftPoint, ok, err := srv.checkLimit(userId, t)
	if err != nil {
		return nil, 0, err
	}
	if !ok {
		return nil, 0, errors.New("已经达到当日最大限制")
	}

	if value >= leftPoint {
		value = leftPoint
	}

	log, err := DefaultPointLimitLogService.UpdateOrCreateLimitLog(UpdateOrCreatePointLimitLogParam{
		Type:            t,
		UserId:          userId,
		AddCurrentValue: value,
		TimePoint:       timeutils.StartOfDay(time.Now()),
	})
	return log, value, err
}
