package point

import (
	"errors"
	"fmt"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/util"
	"time"
)

// 检查是否超过次数
func (c *defaultClientHandle) checkTimes(times int64) error {
	if times == 0 {
		// 等于0表示不限次数
		return nil
	}
	count, err := c.plugin.history.Count(&entity.PointCollectHistory{
		OpenId: c.clientHandle.OpenId,
		Type:   string(c.clientHandle.Type),
		Date:   model.Date{Time: time.Now()},
	})
	if err != nil {
		return err
	}
	if count >= times {
		return errors.New("超过当日次数")
	}
	return nil
}
func (c *defaultClientHandle) checkTimes2(times int64) error {
	if times == 0 {
		// 等于0表示不限次数
		return nil
	}
	//获取次数记录
	result := c.plugin.transactionLimit.FindBy(repository.FindPointTransactionCountLimitBy{
		OpenId:          c.clientHandle.OpenId,
		TransactionType: entity.PointTransactionType(c.clientHandle.Type),
		TransactionDate: model.Date{Time: time.Now()},
	})
	if result.Id == 0 {
		//创建记录
		pointTransActionLimit := entity.PointTransactionCountLimit{
			OpenId:          c.clientHandle.OpenId,
			TransactionType: entity.PointTransactionType(c.clientHandle.Type),
			MaxCount:        entity.PointCollectLimitMap[entity.PointTransactionType(c.clientHandle.Type)],
			CurrentCount:    0,
			UpdateTime:      model.Time{Time: time.Now()},
			TransactionDate: model.Date{Time: time.Now()},
		}
		err := c.saveTransActionLimit(pointTransActionLimit)
		if err != nil {
			return err
		}
		return nil
	}
	if result.CurrentCount >= int(100) {
		return errors.New("超过当日次数")
	}
	//更新记录
	result.CurrentCount += 1
	err := c.updateTransActionLimit(result)
	if err != nil {
		return err
	}
	return nil
}

// 检查积分是否足够
func (c *defaultClientHandle) checkUsrPoints(num int64) error {
	if c.additional.changeType == "dec" && num-c.clientHandle.point <= 0 {
		return errors.New("积分不足")
	}

	return nil
}

// 幂等
func (c *defaultClientHandle) checkIdempotency() error {
	if !util.DefaultLock.Lock(fmt.Sprintf("%s", "collect"+"_"+c.clientHandle.OpenId), time.Second*5) {
		return errors.New("操作频率过快,请稍后再试")
	}
	return nil
}
