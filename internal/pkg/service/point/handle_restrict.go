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
func (c *clientHandle) checkTimes(times int64) error {
	count, err := repository.DefaultPointCollectHistoryRepository.Count(&entity.PointCollectHistory{
		OpenId: c.OpenId,
		Type:   string(c.Type),
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

// 检查积分是否足够
func (c *clientHandle) checkUsrPoints(num int64) error {
	if c.additional.changeType == "dec" && num-c.Point <= 0 {
		return errors.New("积分不足")
	}
	return nil
}

// 幂等
func (c *clientHandle) checkIdempotency() error {
	if !util.DefaultLock.Lock(fmt.Sprintf("%s", string(c.Type)+"_"+c.OpenId), time.Second*5) {
		return errors.New("操作频率过快,请稍后再试")
	}
	return nil
}
