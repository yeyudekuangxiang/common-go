package service

import (
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	repository2 "mio/internal/pkg/repository"
	"mio/pkg/errno"
)

var DefaultPointService = NewPointService(repository2.DefaultPointRepository)

func NewPointService(repo repository2.PointRepository) PointService {
	return PointService{repo: repo}
}

type PointService struct {
	repo repository2.PointRepository
}

// RefreshBalance 根据积分发放记录更新用户剩余积分
func (srv PointService) RefreshBalance(openId string) error {
	transactionList := DefaultPointTransactionService.GetListBy(repository2.GetPointTransactionListBy{
		OpenId: openId,
	})
	balance := 0
	for _, t := range transactionList {
		balance += t.Value
	}

	point := srv.repo.FindBy(repository2.FindPointBy{
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
	return repository2.DefaultPointRepository.Save(&point)
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
		return &entity.Point{}, errno.ErrUserNotFound
	}
	point := srv.repo.FindBy(repository2.FindPointBy{
		OpenId: user.OpenId,
	})
	return &point, nil
}

// FindByOpenId 获取用户积分
func (srv PointService) FindByOpenId(openId string) (*entity.Point, error) {
	if openId == "" {
		return &entity.Point{}, errno.ErrUserNotFound
	}
	point := srv.repo.FindBy(repository2.FindPointBy{
		OpenId: openId,
	})
	return &point, nil
}
