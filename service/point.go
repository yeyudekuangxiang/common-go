package service

import (
	"errors"
	"mio/core/app"
	"mio/model/entity"
	"mio/repository"
)

var DefaultPointService = NewPointService(repository.DefaultPointRepository)

func NewPointService(repo repository.PointRepository) PointService {
	return PointService{repo: repo}
}

type PointService struct {
	repo repository.PointRepository
}

// RefreshBalance 根据积分发放记录更新用户剩余积分
func (srv PointService) RefreshBalance(openId string) error {
	transactionList := DefaultPointTransactionService.GetListBy(repository.GetPointTransactionListBy{
		OpenId: openId,
	})
	balance := 0
	for _, t := range transactionList {
		balance += t.Value
	}

	point := srv.repo.FindBy(repository.FindPointBy{
		OpenId: openId,
	})
	if point.Id == 0 {
		point.OpenId = openId
	}
	if balance < 0 {
		app.Logger.Error("用户积分异常,请检查", openId, balance)
		balance = 0
	}
	point.Balance = balance
	return repository.DefaultPointRepository.Save(&point)
}
func (srv PointService) RefreshBalanceByMq(openId string) {
	err := initUserFlowPool.Submit(func() {
		err := srv.RefreshBalance(openId)
		if err != nil {
			app.Logger.Error("计算用户积分失败:", openId, err)
		}
	})
	if err != nil {
		app.Logger.Error("提交计算用户积分失败:", openId, err)
	}
	return
}

// FindByUserId 获取用户积分
func (srv PointService) FindByUserId(userId int64) (*entity.Point, error) {
	user, err := DefaultUserService.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	if user.OpenId == "" {
		return nil, errors.New("用户不存在")
	}
	point := srv.repo.FindBy(repository.FindPointBy{
		OpenId: user.OpenId,
	})
	return &point, nil
}
