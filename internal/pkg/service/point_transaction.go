package service

import (
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	repository2 "mio/internal/pkg/repository"
	"mio/internal/pkg/util"
	"time"
)

var DefaultPointTransactionService = NewPointTransactionService(repository2.DefaultPointTransactionRepository)

func NewPointTransactionService(repo repository2.PointTransactionRepository) PointTransactionService {
	return PointTransactionService{
		repo: repo,
	}
}

type PointTransactionService struct {
	repo repository2.PointTransactionRepository
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
func (p PointTransactionService) GetListBy(by repository2.GetPointTransactionListBy) []entity.PointTransaction {
	return repository2.DefaultPointTransactionRepository.GetListBy(by)
}
func (p PointTransactionService) FindBy(by repository2.FindPointTransactionBy) (*entity.PointTransaction, error) {
	pt := p.repo.FindBy(by)
	return &pt, nil
}
