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
func (c *defaultClientHandle) checkTimes2() error {
	//获取次数记录
	result := c.plugin.transactionLimit.FindBy(repository.FindPointTransactionCountLimitBy{
		OpenId:          c.clientHandle.OpenId,
		TransactionType: c.clientHandle.Type,
		TransactionDate: model.Date{Time: time.Now()},
	})
	if result.Id == 0 {
		//创建记录
		pointTransActionLimit := entity.PointTransactionCountLimit{
			OpenId:          c.clientHandle.OpenId,
			TransactionType: c.clientHandle.Type,
			MaxCount:        entity.PointCollectLimitMap[c.clientHandle.Type],
			CurrentCount:    1,
			UpdateTime:      model.Time{Time: time.Now()},
			TransactionDate: model.Date{Time: time.Now()},
		}
		err := c.saveTransActionLimit(pointTransActionLimit)
		if err != nil {
			return err
		}
		return nil
	}
	if result.CurrentCount+1 >= result.MaxCount {
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

func (c *defaultClientHandle) checkMaxPoint(maxPoint int64, currPoint int64) error {
	today, _, err := c.plugin.transaction.CountByToday(repository.GetPointTransactionCountBy{
		OpenIds: []string{c.clientHandle.OpenId},
		Type:    c.clientHandle.Type,
	})
	if err != nil {
		return err
	}
	var point int64
	for _, item := range today {
		point += item["value"].(int64)
	}
	if maxPoint-point <= 0 {
		return errors.New("今日积分获取已达到上限")
	}
	if currPoint+point >= maxPoint {
		c.clientHandle.point = maxPoint
	} else {
		c.clientHandle.point = currPoint
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

func (c *defaultClientHandle) checkOrderId(orderId string) error {
	if orderId == "" {
		return errors.New("参数错误:订单号为空")
	}
	result, err := c.plugin.transaction.FindOrder(orderId)
	if err != nil {
		return err
	}
	if result.ID == 0 {
		return nil
	}
	return errors.New("订单号已经存在")
}
