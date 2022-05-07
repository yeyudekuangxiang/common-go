package service

import (
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	duibaApi "mio/pkg/duiba/api/model"
	"mio/pkg/errno"
)

var DefaultDuiBaOrderService = DuiBaOrderService{repo: repository.DefaultDuiBaOrderRepository}

type DuiBaOrderService struct {
	repo repository.DuiBaOrderRepository
}

func (srv DuiBaOrderService) FindByOrderId(orderId string) (*entity.DuiBaOrder, error) {
	order := srv.repo.FindByOrderId(orderId)
	return &order, nil
}

func (srv DuiBaOrderService) CreateOrUpdate(orderId string, info duibaApi.OrderInfo) (*entity.DuiBaOrder, error) {
	user, err := DefaultUserService.GetUserByOpenId(info.Uid)
	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		return nil, errno.ErrUserNotFound
	}

	order := srv.repo.FindByOrderId(orderId)
	if order.ID == 0 {
		order = entity.DuiBaOrder{
			OrderNum:         info.OrderNum,
			DevelopBizId:     info.DevelopBizId,
			CreateTime:       info.CreateTime.ToInt(),
			FinishTime:       info.FinishTime.ToInt(),
			TotalCredits:     int(info.TotalCredits.ToInt()),
			ConsumerPayPrice: info.ConsumerPayPrice.ToFloat(),
			Source:           info.Source,
			OrderStatus:      info.OrderStatus,
			ErrorMsg:         info.ErrorMsg,
			Type:             info.Type,
			ExpressPrice:     info.ExpressPrice,
			Account:          info.Account,
			OrderItemList:    string(info.OrderItemList),
			ReceiveAddrInfo:  string(info.ReceiveAddrInfo),
			UserId:           user.ID,
			OrderId:          orderId,
		}
		return &order, srv.repo.Create(&order)
	}

	order.FinishTime = info.FinishTime.ToInt()
	order.OrderStatus = info.OrderStatus
	order.ErrorMsg = info.ErrorMsg
	order.OrderItemList = string(info.OrderItemList)
	return &order, srv.repo.Save(&order)

}
