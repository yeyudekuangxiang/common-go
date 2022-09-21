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
	var limitN, i, it int
	var ok bool
	var msg string

	if i, ok = entity.PointCollectLimitMap[t]; ok {
		limitN = i
		msg = "今日"
	} else if i, ok = entity.PointCollectLimitOnceMap[t]; ok {
		it = 1
		limitN = i
	}

	if !ok {
		return nil
	}

	limitWhere := repository.FindPointTransactionCountLimitBy{
		OpenId:          openId,
		TransactionType: t,
	}
	if it == 0 {
		limitWhere.TransactionDate = model.Date{Time: time.Now()}
	}

	limit := srv.repo.FindBy(limitWhere)

	if limit.Id == 0 {
		newLimit, err := srv.createLimitOfToday(t, openId)
		if err != nil {
			return err
		}
		limit = *newLimit
	}

	if limit.CurrentCount >= limitN {
		return errors.New("获取积分次数达到" + msg + "上限")
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
	var maxCount int
	if max, ok := entity.PointCollectLimitMap[transactionType]; ok {
		maxCount = max
	} else if max, ok = entity.PointCollectLimitOnceMap[transactionType]; ok {
		maxCount = max
	}
	limit := entity.PointTransactionCountLimit{
		OpenId:          openId,
		TransactionType: transactionType,
		MaxCount:        maxCount,
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

//每天最多获取多少分 currPoint : 要加的积分。 maxPoint : 上限积分
func (srv PointTransactionCountLimitService) CheckMaxPoint(transactionType entity.PointTransactionType, openId string, currPoint int64, maxPoint int64) (int64, error) {
	today, _, err := NewPointTransactionService(srv.ctx).CountByToday(openId, transactionType)
	if err != nil {
		return 0, err
	}
	var todayValue int64
	for _, item := range today {
		todayValue += item["value"].(int64)
	}
	if maxPoint-todayValue <= 0 {
		return 0, errors.New("今日积分已达上限")
	}
	if maxPoint-todayValue > 0 && maxPoint-todayValue < currPoint {
		return maxPoint - todayValue, nil
	}
	//正常加积分
	return currPoint, nil
}

func (srv PointTransactionCountLimitService) CheckMaxPointByMonth(transactionType entity.PointTransactionType, openId string, currPoint int64, maxPoint int64) (int64, error) {
	month, _, err := NewPointTransactionService(srv.ctx).CountByMonth(openId, transactionType)
	if err != nil {
		return 0, err
	}
	var todayValue int64
	for _, item := range month {
		todayValue += item["value"].(int64)
	}
	if maxPoint-todayValue <= 0 {
		return 0, errors.New("此回收分类已达到本月获取积分上限")
	}
	if maxPoint-todayValue > 0 && maxPoint-todayValue < currPoint {
		return maxPoint - todayValue, nil
	}
	//正常加积分
	return currPoint, nil
}
