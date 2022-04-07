package service

import (
	"mio/internal/util"
	"mio/model"
	"mio/model/entity"
	"mio/repository"
	"time"
)

var DefaultPointTransactionService = NewPointTransactionService(repository.DefaultPointTransactionRepository)

func NewPointTransactionService(repo repository.PointTransactionRepository) PointTransactionService {
	return PointTransactionService{
		repo: repo,
	}
}

type PointTransactionService struct {
	repo repository.PointTransactionRepository
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

	err = p.repo.Save(&transaction)
	if err != nil {
		return nil, err
	}

	_ = DefaultPointService.RefreshBalance(param.OpenId)

	return &transaction, nil
}

// GetListBy 查询记录列表
func (p PointTransactionService) GetListBy(by repository.GetPointTransactionListBy) []entity.PointTransaction {
	return repository.DefaultPointTransactionRepository.GetListBy(by)
}
func (p PointTransactionService) FindBy(by repository.FindPointTransactionBy) (*entity.PointTransaction, error) {
	pt := p.repo.FindBy(by)
	return &pt, nil
}
