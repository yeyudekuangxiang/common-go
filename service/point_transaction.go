package service

import (
	"mio/internal/util"
	"mio/model"
	"mio/model/entity"
	"mio/repository"
	"time"
)

var DefaultPointTransactionService = PointTransactionService{}

type PointTransactionService struct {
}

// Create 添加发放积分记录并且更新用户剩余积分
func (p PointTransactionService) Create(param CreatePointTransactionParam) (*entity.PointTransaction, error) {

	err := DefaultPointTransactionCountLimitService.CheckLimitAndUpdate(param.Type, param.OpenId)
	if err != nil {
		return nil, err
	}

	transaction := entity.PointTransaction{
		OpenId:         param.OpenId,
		TransactionId:  util.UUID(),
		Type:           param.Type,
		Value:          param.Value,
		CreateTime:     model.Time{Time: time.Now()},
		AdditionalInfo: param.AdditionInfo,
	}

	err = repository.DefaultPointTransactionRepository.Save(&transaction)
	if err != nil {
		return nil, err
	}

	err = DefaultPointService.RefreshBalance(param.OpenId)

	return &transaction, err
}

// GetListBy 查询记录列表
func (p PointTransactionService) GetListBy(by repository.GetPointTransactionListBy) []entity.PointTransaction {
	return repository.DefaultPointTransactionRepository.GetListBy(by)
}
