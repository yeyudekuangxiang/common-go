package service

import (
	"errors"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"time"
)

type PointTransactionCountLimitService struct {
	ctx  *context.MioContext
	repo *repository.PointTransactionCountLimitRepository
}

func NewPointTransactionCountLimitService(ctx *context.MioContext) *PointTransactionCountLimitService {
	return &PointTransactionCountLimitService{ctx: ctx, repo: repository.NewPointTransactionCountLimitRepository(ctx)}
}

// CheckLimitAndUpdate 检查积分发送次数限制
func (srv PointTransactionCountLimitService) CheckLimitAndUpdate(t entity.PointTransactionType, openId string) error {
	limitNum, ok := entity.PointCollectLimitMap[t]
	if !ok {
		return nil
	}
	limit := srv.repo.FindBy(repository.FindPointTransactionCountLimitBy{
		OpenId:          openId,
		TransactionType: t,
		TransactionDate: model.Date{Time: time.Now()},
	})

	if limit.Id == 0 {
		newLimit, err := srv.createLimitOfToday(t, openId)
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
func (srv PointTransactionCountLimitService) createLimitOfToday(transactionType entity.PointTransactionType, openId string) (*entity.PointTransactionCountLimit, error) {
	limit := entity.PointTransactionCountLimit{
		OpenId:          openId,
		TransactionType: transactionType,
		MaxCount:        entity.PointCollectLimitMap[transactionType],
		CurrentCount:    0,
		UpdateTime:      model.Time{Time: time.Now()},
		TransactionDate: model.Date{Time: time.Now()},
	}
	return &limit, srv.repo.Save(&limit)
}
func (srv PointTransactionCountLimitService) CheckLimit(transactionType entity.PointTransactionType, openId string) (bool, error) {
	limitNum, ok := entity.PointCollectLimitMap[transactionType]
	if !ok {
		return true, nil
	}

	limit := srv.repo.FindBy(repository.FindPointTransactionCountLimitBy{
		OpenId:          openId,
		TransactionType: transactionType,
		TransactionDate: model.Date{Time: time.Now()},
	})
	if limit.Id == 0 {
		newLimit, err := srv.createLimitOfToday(transactionType, openId)
		if err != nil {
			return false, err
		}
		limit = *newLimit
	}

	return limit.CurrentCount < limitNum, nil
}

//每天最多获取多少分
func (srv PointTransactionCountLimitService) CheckMaxPoint(transactionType entity.PointTransactionType, openId string, currPoint *int64) error {
	dayPoint, ok := entity.PointCollectValueMapDay[transactionType]
	if !ok {
		return errors.New("未匹配到日积分限制规则")
	}
	today, _, err := NewPointTransactionService(srv.ctx).CountByToday([]string{openId}, transactionType)
	if err != nil {
		return err
	}
	var todayValue int64
	for _, item := range today {
		todayValue += item["value"].(int64)
	}
	point := int64(dayPoint)
	if point-todayValue <= 0 {
		return errors.New("今日积分已达上限")
	}
	if point-todayValue > 0 && point-todayValue < *currPoint {
		*currPoint = point - todayValue
	}
	//正常加积分
	return nil
}
