package service

import (
	"errors"
	"mio/model"
	"mio/model/entity"
	"mio/repository"
	"time"
)

var DefaultPointTransactionCountLimitService = PointTransactionCountLimitService{}

type PointTransactionCountLimitService struct {
}

// CheckLimitAndUpdate 检查积分发送次数限制
func (p PointTransactionCountLimitService) CheckLimitAndUpdate(t entity.PointTransactionType, openId string) error {
	limitNum, ok := entity.PointCollectLimitMap[t]
	if !ok {
		return nil
	}
	limit := repository.DefaultPointTransactionCountLimitRepository.FindBy(repository.FindPointTransactionCountLimitBy{
		OpenId:          openId,
		TransactionType: t,
		TransactionDate: model.Date{Time: time.Now()},
	})
	if limit.Id == 0 {
		_, err := p.createLimitOfToday(t, openId)
		if err != nil {
			return err
		}
		return nil
	}

	if limit.CurrentCount >= limitNum {
		return errors.New("达到当日该类别最大积分限制")
	}

	limit.MaxCount++
	err := repository.DefaultPointTransactionCountLimitRepository.Save(&limit)
	if err != nil {
		return err
	}
	return nil
}

//创建用户今天积分发送次数限制记录
func (p PointTransactionCountLimitService) createLimitOfToday(transactionType entity.PointTransactionType, openId string) (*entity.PointTransactionCountLimit, error) {
	limit := entity.PointTransactionCountLimit{
		OpenId:          openId,
		TransactionType: transactionType,
		MaxCount:        entity.PointCollectLimitMap[transactionType],
		CurrentCount:    1,
		UpdateTime:      model.Time{Time: time.Now()},
		TransactionDate: model.Date{Time: time.Now()},
	}
	return &limit, repository.DefaultPointTransactionCountLimitRepository.Save(&limit)
}
