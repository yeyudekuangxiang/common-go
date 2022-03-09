package service

import "mio/repository"

var DefaultPointService = PointService{}

type PointService struct {
}

// RefreshBalance 根据积分发放记录更新用户剩余积分
func (p PointService) RefreshBalance(openId string) error {
	transactionList := DefaultPointTransactionService.GetListBy(repository.GetPointTransactionListBy{
		OpenId: openId,
	})
	balance := 0
	for _, t := range transactionList {
		balance += t.Value
	}

	point := repository.DefaultPointRepository.FindBy(repository.FindPointBy{
		OpenId: openId,
	})
	if point.Id == 0 {
		point.OpenId = openId
	}
	point.Balance = balance
	return repository.DefaultPointRepository.Save(&point)
}
