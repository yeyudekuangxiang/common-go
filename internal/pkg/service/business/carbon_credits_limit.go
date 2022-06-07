package business

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	ebusiness "mio/internal/pkg/model/entity/business"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/timeutils"
	"time"
)

var DefaultCarbonCreditsLimitService = CarbonCreditsLimitService{}

type CarbonCreditsLimitService struct {
}

//检测是否超出限制 返回还可领取次数 还可领取积分值 是否超出
func (srv CarbonCreditsLimitService) checkLimit(userId int64, t ebusiness.CarbonType) (int, decimal.Decimal, bool, error) {

	//需要方法 根据用户id查询用户信息
	userInfo := ebusiness.User{}
	carbonScene, err := DefaultCarbonSceneService.FindScene(t)
	if err != nil {
		return 0, decimal.Decimal{}, false, err
	}
	if carbonScene.ID == 0 {
		return 0, decimal.Decimal{}, false, errors.New("未查询到此低碳场景")
	}
	companyCarbonScene, err := DefaultCompanyCarbonSceneService.FindCompanyScene(FindCompanyCarbonSceneParam{
		CompanyId: userInfo.BCompanyId,
	})
	if err != nil {
		return 0, decimal.Decimal{}, false, err
	}
	if companyCarbonScene.ID == 0 {
		return 0, decimal.Decimal{}, false, errors.New("未查询到此低碳场景")
	}

	limitLog, err := DefaultCarbonCreditsLimitLogService.FindLimitLog(FindCarbonCreditsLimitLogParam{
		TimePoint: timeutils.Now().StartOfDay().Time,
		Type:      t,
		UserId:    userId,
	})
	if err != nil {
		return 0, decimal.Decimal{}, false, err
	}

	if limitLog.ID == 0 {
		return companyCarbonScene.MaxCount, carbonScene.MaxCarbonCredits, true, nil
	}

	count := companyCarbonScene.MaxCount - limitLog.CurrentCount
	if count <= 0 {
		return 0, decimal.Decimal{}, false, nil
	}

	credits := carbonScene.MaxCarbonCredits.Sub(limitLog.CurrentValue)
	if credits.LessThanOrEqual(decimal.Zero) {
		return 0, decimal.Decimal{}, false, nil
	}

	return count, credits, true, nil
}

// CheckLimit 加读写锁检测是否超出限制 返回还可领取次数 还可领取积分值 是否超出
func (srv CarbonCreditsLimitService) CheckLimit(userId int64, t ebusiness.CarbonType) (int, decimal.Decimal, bool, error) {
	lockKey := fmt.Sprintf("CheckCarbonLimit%d%s", userId, t)
	util.DefaultLock.LockWait(lockKey, time.Second*5)
	defer util.DefaultLock.UnLock(lockKey)
	return srv.checkLimit(userId, t)
}

//CheckLimitAndUpdate 加读写锁检测是否超出限制并且修改值 返回记录 以及实际新增的积分数量
func (srv CarbonCreditsLimitService) CheckLimitAndUpdate(userId int64, value decimal.Decimal, t ebusiness.CarbonType) (*ebusiness.CarbonCreditsLimitLog, decimal.Decimal, error) {
	lockKey := fmt.Sprintf("CheckCarbonLimit%d%s", userId, t)
	util.DefaultLock.LockWait(lockKey, time.Second*5)
	defer util.DefaultLock.UnLock(lockKey)
	_, credits, ok, err := srv.checkLimit(userId, t)
	if err != nil {
		return nil, decimal.Decimal{}, err
	}
	if !ok {
		return nil, decimal.Decimal{}, errors.New("已经达到当日最大限制")
	}

	if value.GreaterThan(credits) {
		value = credits
	}

	log, err := DefaultCarbonCreditsLimitLogService.UpdateOrCreateLimitLog(UpdateOrCreateCarbonCreditsLimitLogParam{
		Type:            t,
		UserId:          userId,
		AddCurrentValue: value,
		TimePoint:       timeutils.StartOfDay(time.Now()),
	})
	return log, value, err
}
