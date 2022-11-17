package service

import (
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service/track"
	"mio/internal/pkg/util/timeutils"
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

	//只有成功才能上报到诸葛
	if info.OrderStatus == "success" {
		err = srv.OrderMaiDian(info, user.ID, user.OpenId)
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

type OrderItem struct {
	Title     string `json:"title"`
	IsSelf    string `json:"isSelf"`
	PerCredit string `json:"perCredit"`
	PerPrice  string `json:"perPrice"`
}

func (srv DuiBaOrderService) OrderMaiDian(order duibaApi.OrderInfo, uid int64, openid string) error {
	//redisKey := "mp2c:order_duiba_to_zhuge_test_v6"
	typeName := order.Type
	switch order.Type {
	case "coupon":
		typeName = "虚拟券"
		break
	case "virtual":
		typeName = "充值商品"
		break
	case "object":
		typeName = "实物商品"
		break
	}
	statusName := order.OrderStatus
	switch order.OrderStatus {
	case "fail":
		statusName = "失败"
		break
	case "afterSend":
		statusName = "已发货"
		break
	case "success":
		statusName = "成功"
		break
	case "waitSend":
		statusName = "待发货"
		break
	}
	//上报到诸葛
	zhuGeAttr := make(map[string]interface{}, 0)
	zhuGeAttr["用户uid"] = uid
	zhuGeAttr["兑吧订单号"] = order.OrderNum
	zhuGeAttr["下单时间"] = timeutils.UnixMilli(order.CreateTime.ToInt()).Format(timeutils.TimeFormat)
	zhuGeAttr["消耗积分数"] = order.TotalCredits
	zhuGeAttr["支付金额"] = order.ConsumerPayPrice
	zhuGeAttr["兑换类型"] = order.Source
	zhuGeAttr["商品类型"] = typeName
	zhuGeAttr["运费"] = order.ExpressPrice
	zhuGeAttr["订单状态"] = statusName
	zhuGeAttr["用户openid"] = openid
	//var orderItemList []OrderItem
	//json.Unmarshal([]byte(order.OrderItemList), &orderItemList)
	for _, item := range order.OrderItemList.OrderItemList() {
		orderItemType := ""
		switch item.IsSelf {
		case "1":
			orderItemType = "自有"
			break
		case "0":
			orderItemType = "兑吧"
			break
		}
		zhuGeAttr["商品名称"] = item.Title
		zhuGeAttr["商品来源"] = orderItemType
		zhuGeAttr["所需兑换积分"] = item.PerCredit
		zhuGeAttr["所需兑换金额"] = item.PerPrice
		break
	}
	duibaOrder := srv.repo.FindByUid(uid)
	if duibaOrder.ID == 0 {
		zhuGeAttr["是否首单"] = "是"
	} else {
		zhuGeAttr["是否首单"] = "否"
	}
	app.Logger.Infof("商场订单，积分打点失败 %+v %v %+v", config.ZhuGeEventName.DuiBaOrder, openid, zhuGeAttr)

	track.DefaultZhuGeService().Track(config.ZhuGeEventName.DuiBaOrder, openid, zhuGeAttr)
	return nil
}
