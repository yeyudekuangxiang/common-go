package business

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	ebusiness "mio/internal/pkg/model/entity/business"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/timeutils"
	"mio/pkg/errno"
	"time"
)

var DefaultCarbonCreditsLimitService = CarbonCreditsLimitService{}

type CarbonCreditsLimitService struct {
}

//检测是否超出限制 返回还可领取次数
func (srv CarbonCreditsLimitService) checkLimit(userId int64, t ebusiness.CarbonType) (int, error) {

	userInfo, err := DefaultUserService.GetBusinessUserById(userId)
	if err != nil {
		return 0, err
	}
	if userInfo.ID == 0 {
		return 0, errno.ErrUserNotFound
	}

	carbonScene, err := DefaultCarbonSceneService.FindScene(t)
	if err != nil {
		return 0, err
	}
	if carbonScene.ID == 0 {
		return 0, errors.New("未查询到此低碳场景")
	}
	companyCarbonScene, err := DefaultCompanyCarbonSceneService.FindCompanyScene(FindCompanyCarbonSceneParam{
		CompanyId:     userInfo.BCompanyId,
		CarbonSceneId: carbonScene.ID,
	})
	if err != nil {
		return 0, err
	}
	if companyCarbonScene.ID == 0 {
		return 0, errors.New("未查询到此低碳场景")
	}

	limitLog, err := DefaultCarbonCreditsLimitLogService.FindLimitLog(FindCarbonCreditsLimitLogParam{
		TimePoint: timeutils.Now().StartOfDay().Time(),
		Type:      t,
		UserId:    userId,
	})
	if err != nil {
		return 0, err
	}

	if limitLog.ID == 0 {
		return companyCarbonScene.MaxCount, nil
	}

	count := companyCarbonScene.MaxCount - limitLog.CurrentCount
	if count <= 0 {
		return 0, nil
	}

	return count, nil
}

// CheckLimit 加读写锁检测是否超出限制 返回还可领取次数 还可领取积分值 是否超出
func (srv CarbonCreditsLimitService) CheckLimit(userId int64, t ebusiness.CarbonType) (int, error) {
	lockKey := fmt.Sprintf("CheckCarbonLimit%d%s", userId, t)
	util.DefaultLock.LockWait(lockKey, time.Second*5)
	defer util.DefaultLock.UnLock(lockKey)
	return srv.checkLimit(userId, t)
}

//CheckLimitAndUpdate 加读写锁检测是否超出限制并且修改值 返回记录
func (srv CarbonCreditsLimitService) CheckLimitAndUpdate(userId int64, value decimal.Decimal, t ebusiness.CarbonType) (*ebusiness.CarbonCreditsLimitLog, error) {
	lockKey := fmt.Sprintf("CheckCarbonLimit%d%s", userId, t)
	util.DefaultLock.LockWait(lockKey, time.Second*5)
	defer util.DefaultLock.UnLock(lockKey)
	count, err := srv.checkLimit(userId, t)
	if err != nil {
		return nil, err
	}
	if count <= 0 {
		return nil, errors.New("已经达到当日最大限制")
	}

	log, err := DefaultCarbonCreditsLimitLogService.UpdateOrCreateLimitLog(UpdateOrCreateCarbonCreditsLimitLogParam{
		Type:            t,
		UserId:          userId,
		AddCurrentValue: value,
		TimePoint:       timeutils.StartOfDay(time.Now()),
	})
	return log, err
}
