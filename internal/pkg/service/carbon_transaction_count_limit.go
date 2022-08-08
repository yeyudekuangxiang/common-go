package service

import (
	"errors"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/repository/repotypes"
	"time"
)

type CarbonTransactionCountLimitService struct {
	ctx  *context.MioContext
	repo repository.CarbonTransactionCountLimitRepository
}

func NewCarbonTransactionCountLimitService(ctx *context.MioContext) *CarbonTransactionCountLimitService {
	return &CarbonTransactionCountLimitService{ctx: ctx, repo: repository.NewCarbonTransactionCountLimitRepository(ctx)}
}

//  检查积分发送次数限制

func (srv CarbonTransactionCountLimitService) CheckLimitAndUpdate(carbonTransactionType entity.CarbonTransactionType, openId string) error {
	limitNum, ok := entity.CarbonCollectLimitMap[carbonTransactionType]
	if !ok {
		return nil
	}
	limit := srv.repo.FindBy(repotypes.FindCarbonTransactionCountLimitFindByDO{
		OpenId: openId,
		Type:   carbonTransactionType,
		VDate:  time.Now(),
	})
	if limit.ID != 0 {
		newLimit, err := srv.createLimitOfToday(carbonTransactionType, openId)
		if err != nil {
			return err
		}
		limit = *newLimit
	}
	if limit.CurrentCount >= limitNum {
		return errors.New("达到当日该类别最大积分限制")
	}
	limit.CurrentCount++
	err := srv.repo.Save(&limit)
	if err != nil {
		return err
	}
	return nil
}

//创建用户今天积分发送次数限制记录
func (srv CarbonTransactionCountLimitService) createLimitOfToday(transactionType entity.CarbonTransactionType, openId string) (*entity.CarbonTransactionCountLimitDay, error) {
	limit := entity.CarbonTransactionCountLimitDay{
		OpenId:       openId,
		Type:         transactionType,
		MaxCount:     entity.CarbonCollectLimitMap[transactionType],
		CurrentCount: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		VDate:        time.Now(),
	}
	return &limit, srv.repo.Save(&limit)
}

func (srv CarbonTransactionCountLimitService) CheckLimit(transactionType entity.CarbonTransactionType, openId string) (bool, error) {
	limitNum, ok := entity.CarbonCollectLimitMap[transactionType]
	if !ok {
		return true, nil
	}
	limit := srv.repo.FindBy(repotypes.FindCarbonTransactionCountLimitFindByDO{
		OpenId: openId,
		Type:   transactionType,
		VDate:  time.Now(),
	})
	if limit.ID == 0 {
		newLimit, err := srv.createLimitOfToday(transactionType, openId)
		if err != nil {
			return false, err
		}
		limit = *newLimit
	}

	return limit.CurrentCount < limitNum, nil
}
