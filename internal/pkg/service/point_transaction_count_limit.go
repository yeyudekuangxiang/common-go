package service

import (
	"errors"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	repository2 "mio/internal/pkg/repository"
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
	limit := repository2.DefaultPointTransactionCountLimitRepository.FindBy(repository2.FindPointTransactionCountLimitBy{
		OpenId:          openId,
		TransactionType: t,
		TransactionDate: model.Date{Time: time.Now()},
	})

	if limit.Id == 0 {
		newLimit, err := p.createLimitOfToday(t, openId)
		if err != nil {
			return err
		}
		limit = *newLimit
	}

	if limit.CurrentCount >= limitNum {
		return errors.New("达到当日该类别最大积分限制")
	}

	limit.CurrentCount++
	err := repository2.DefaultPointTransactionCountLimitRepository.Save(&limit)
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
		CurrentCount:    0,
		UpdateTime:      model.Time{Time: time.Now()},
		TransactionDate: model.Date{Time: time.Now()},
	}
	return &limit, repository2.DefaultPointTransactionCountLimitRepository.Save(&limit)
}
func (p PointTransactionCountLimitService) CheckLimit(transactionType entity.PointTransactionType, openId string) (bool, error) {
	limitNum, ok := entity.PointCollectLimitMap[transactionType]
	if !ok {
		return true, nil
	}

	limit := repository2.DefaultPointTransactionCountLimitRepository.FindBy(repository2.FindPointTransactionCountLimitBy{
		OpenId:          openId,
		TransactionType: transactionType,
		TransactionDate: model.Date{Time: time.Now()},
	})
	if limit.Id == 0 {
		newLimit, err := p.createLimitOfToday(transactionType, openId)
		if err != nil {
			return false, err
		}
		limit = *newLimit
	}

	return limit.CurrentCount < limitNum, nil
}
