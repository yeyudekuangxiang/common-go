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

//CheckLimitAndUpdate 检查积分发送次数限制
func (srv CarbonTransactionCountLimitService) CheckLimitAndUpdate(carbonTransactionType entity.CarbonTransactionType, openId string, limitNum int) error {
	limit := srv.repo.FindBy(repotypes.FindCarbonTransactionCountLimitFindByDO{
		OpenId: openId,
		Type:   carbonTransactionType,
		VDate:  time.Now().Format("2006-01-02"),
	})
	if limit.ID == 0 {
		newLimit, err := srv.createLimitOfToday(carbonTransactionType, openId, limitNum)
		if err != nil {
			return err
		}
		limit = *newLimit
	}
	if limit.CurrentCount >= limitNum {
		return errors.New("达到当日该类别最大碳量限制")
	}
	limit.CurrentCount++
	err := srv.repo.Save(&limit)
	if err != nil {
		return err
	}
	return nil
}

//创建用户今天积分发送次数限制记录
func (srv CarbonTransactionCountLimitService) createLimitOfToday(transactionType entity.CarbonTransactionType, openId string, limitNum int) (*entity.CarbonTransactionCountLimit, error) {
	limit := entity.CarbonTransactionCountLimit{
		OpenId:       openId,
		Type:         transactionType,
		MaxCount:     limitNum,
		CurrentCount: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		VDate:        time.Now().Format("2006-01-02"),
	}
	return &limit, srv.repo.Save(&limit)
}

func (srv CarbonTransactionCountLimitService) CheckLimit(transactionType entity.CarbonTransactionType, openId string, limitNum int) (bool, error) {
	limit := srv.repo.FindBy(repotypes.FindCarbonTransactionCountLimitFindByDO{
		OpenId: openId,
		Type:   transactionType,
		VDate:  time.Now().Format("2006-01-02"),
	})
	if limit.ID == 0 {
		newLimit, err := srv.createLimitOfToday(transactionType, openId, limitNum)
		if err != nil {
			return false, err
		}
		limit = *newLimit
	}
	return limit.CurrentCount < limitNum, nil
}
