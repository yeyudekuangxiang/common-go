package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"math/rand"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	repository2 "mio/internal/pkg/repository"
	"mio/internal/pkg/repository/repotypes"
	"mio/internal/pkg/service/event"
	"mio/internal/pkg/service/product"
	"mio/internal/pkg/service/srv_types"
	util2 "mio/internal/pkg/util"
	duibaApi "mio/pkg/duiba/api/model"
	"mio/pkg/duiba/util"
	"mio/pkg/errno"
	"strconv"
	"time"
)

var DefaultOrderService = NewOrderService(repository2.DefaultOrderRepository)

func NewOrderService(repo repository2.OrderRepository) OrderService {
	return OrderService{repo: repo}
}

type OrderService struct {
	repo repository2.OrderRepository
}

// CalculateAndCheck 计算商品价格并且检查库存
func (srv OrderService) CalculateAndCheck(items []repository2.CheckStockItem) (*CalculateProductResult, error) {
	itemIds := make([]string, 0)
	itemMap := make(map[string]repository2.CheckStockItem)
	for _, item := range items {
		itemIds = append(itemIds, item.ItemId)
		itemMap[item.ItemId] = item
	}

	productItems := product.DefaultProductItemService.GetListBy(product.GetProductItemListParam{
		ItemIds: itemIds,
	})
	if len(productItems) != len(itemIds) {
		return nil, errors.New("存在失效商品,请去掉失效商品后重试")
	}

	result := &CalculateProductResult{ItemList: make([]submitOrderItem, 0)}
	for _, productItem := range productItems {
		wantCount := itemMap[productItem.ProductItemId].Count
		if !productItem.Active {
			return nil, errors.New("商品`" + productItem.Title + "`已下架")
		}
		if productItem.RemainingCount < wantCount {
			return nil, errors.New("商品`" + productItem.Title + "`库存不足")
		}
		result.TotalCost += productItem.Cost
		result.ItemList = append(result.ItemList, submitOrderItem{
			ItemId: productItem.ProductItemId,
			Cost:   productItem.Cost,
			Count:  wantCount,
		})
	}

	return result, nil
}

// SubmitOrder 用于外部下单
/*
原java下单流程 目前没有做保存活动记录 商品购买上限 证书发放三个操作
1 生成订单id
2 计算价格
3 保存order_item
4 减库存
5 根据productItemId 查询相关活动并且保存活动记录
6 检查是否达到商品购买上限
7 未达到上限时保存记录 如果Effective_data为空 则根据openid itemid添加或者更新数据 如果Effective_data不为空 则根据当天日期 openid itemid添加或者更新数据
8 根据itemid发放证书
9 创建订单
*/
func (srv OrderService) SubmitOrder(param SubmitOrderParam) (*entity.Order, error) {
	lockKey := fmt.Sprintf("SubmitOrder%d", param.Order.UserId)
	if !util2.DefaultLock.Lock(lockKey, time.Second*10) {
		return nil, errno.ErrLimit.WithCaller()
	}
	defer util2.DefaultLock.UnLock(lockKey)

	checkItems := make([]repository2.CheckStockItem, 0)
	for _, item := range param.Items {
		checkItems = append(checkItems, repository2.CheckStockItem{
			ItemId: item.ItemId,
			Count:  item.Count,
		})
	}
	calculateResult, err := srv.CalculateAndCheck(checkItems)
	if err != nil {
		return nil, err
	}

	return srv.submitOrder(submitOrderParam{
		Order: submitOrder{
			AddressId: param.Order.AddressId,
			UserId:    param.Order.UserId,
			TotalCost: calculateResult.TotalCost,
			OrderType: entity.OrderTypePurchase,
		},
		Items:           calculateResult.ItemList,
		PartnershipType: param.PartnershipType,
	})
}

//submitOrder (此方法可自定义需要支付的金额 需谨慎使用 用户下单请使用 OrderService.SubmitOrder 方法创建订单)
func (srv OrderService) submitOrder(param submitOrderParam) (*entity.Order, error) {
	//防止并发
	if !util2.DefaultLock.Lock("submitOrder_"+strconv.FormatInt(param.Order.UserId, 10), time.Second*5) {
		return nil, errors.New("操作频率过快,请稍后再试")
	}

	orderSuccess := false
	orderId := util2.UUID()

	//查询用户信息
	user, err := DefaultUserService.GetUserById(param.Order.UserId)
	if err != nil {
		return nil, err
	}
	if user.ID == 0 || user.OpenId == "" {
		return nil, errors.New("未查找到用户信息,请联系管理员")
	}

	//检查积分
	point, err := DefaultPointService.FindByUserId(param.Order.UserId)
	if err != nil {
		return nil, err
	}
	if point.Balance < param.Order.TotalCost {
		return nil, errors.New("积分不足,无法兑换")
	}

	//检查并且锁定库存
	checkStockItems := make([]repository2.CheckStockItem, 0)
	for _, item := range param.Items {
		checkStockItems = append(checkStockItems, repository2.CheckStockItem{
			ItemId: item.ItemId,
			Count:  item.Count,
		})
	}
	err = product.DefaultProductItemService.CheckAndLockStock(checkStockItems)
	if err != nil {
		return nil, err
	}
	//下单失败释放库存
	defer func() {
		if orderSuccess == false {
			//释放库存
			err := product.DefaultProductItemService.UnLockStock(checkStockItems)
			if err != nil {
				app.Logger.Errorf("释放库存失败 %+v %+v", checkStockItems, err)
			}
		}
	}()

	//扣除积分
	_, err = DefaultPointTransactionService.Create(CreatePointTransactionParam{
		OpenId:       user.OpenId,
		Value:        -param.Order.TotalCost,
		Type:         entity.POINT_PURCHASE,
		AdditionInfo: `{"orderId":"` + orderId + `"}`,
	})
	if err != nil {
		return nil, err
	}
	//下单失败返还积分
	defer func() {
		if !orderSuccess {
			//返还积分
			_, err = DefaultPointTransactionService.Create(CreatePointTransactionParam{
				OpenId:       user.OpenId,
				Value:        param.Order.TotalCost,
				Type:         entity.POINT_ADJUSTMENT,
				AdditionInfo: `{"orderId":"` + orderId + `","message":"下单失败返还积分"}`,
			})
		}
	}()

	//创建订单
	order, orderItems, err := srv.create(orderId, param)
	if err != nil {
		return nil, err
	}

	srv.afterCreateOrder(param, user, order, orderItems)

	orderSuccess = true
	return order, nil
}

//直接创建订单 不会扣除积分(请勿使用此方法创建订单 请使用 OrderService.SubmitOrder 方法创建订单)
func (srv OrderService) create(orderId string, param submitOrderParam) (*entity.Order, []entity.OrderItem, error) {
	var addressId *string
	if param.Order.AddressId != "" {
		addressId = &param.Order.AddressId
	}
	user, err := DefaultUserService.GetUserById(param.Order.UserId)
	if err != nil {
		return nil, nil, err
	}
	if user.ID == 0 || user.OpenId == "" {
		return nil, nil, errors.New("未查找到用户信息,请联系管理员")
	}

	order := &entity.Order{
		OrderId:          orderId,
		AddressId:        addressId,
		OpenId:           user.OpenId,
		TotalCost:        param.Order.TotalCost,
		Status:           entity.OrderStatusPaid,
		PaidTime:         model.NewTime(),
		OrderReferenceId: model.NullString(fmt.Sprintf("%d%d", time.Now().Unix(), int(rand.Float64()*10000))),
		OrderType:        param.Order.OrderType,
	}

	orderItems := make([]entity.OrderItem, 0)
	for _, item := range param.Items {
		orderItems = append(orderItems, entity.OrderItem{
			OrderId: orderId,
			ItemId:  item.ItemId,
			Cost:    item.Cost,
			Count:   item.Count,
		})
	}

	err = srv.repo.SubmitOrder(order, &orderItems)
	if err != nil {
		return nil, nil, err
	}

	return order, orderItems, nil
}

// SubmitOrderForGreenMonday 用于greenmonday活动用户下单
func (srv OrderService) SubmitOrderForGreenMonday(param SubmitOrderForGreenParam) (*entity.Order, error) {
	return srv.submitOrder(submitOrderParam{
		Order: submitOrder{
			UserId:    param.UserId,
			AddressId: param.AddressId,
			OrderType: entity.OrderTypePurchase,
			TotalCost: 1,
		},
		Items: []submitOrderItem{
			{
				ItemId: param.ItemId,
				Count:  1,
				Cost:   1,
			},
		},
	})
}

var duiBaOrderStatusMap = map[duibaApi.OrderStatus]entity.OrderStatus{
	duibaApi.OrderStatusWaitAudit: entity.OrderStatusPaid,
	duibaApi.OrderStatusWaitSend:  entity.OrderStatusPaid,
	duibaApi.OrderStatusAfterSend: entity.OrderStatusInTransit,
	duibaApi.OrderStatusSuccess:   entity.OrderStatusComplete,
	duibaApi.OrderStatusFail:      entity.OrderStatusError,
}

func (srv OrderService) CreateOrUpdateOrderOfDuiBa(orderId string, info duibaApi.OrderInfo) (*entity.Order, error) {
	order := srv.repo.FindByOrderId(orderId)
	if order.ID == 0 {
		return srv.CreateOrderOfDuiBa(orderId, info)
	}
	return srv.UpdateOrderOfDuiBa(orderId, info)
}
func (srv OrderService) UpdateOrderOfDuiBa(orderId string, info duibaApi.OrderInfo) (*entity.Order, error) {
	order := srv.repo.FindByOrderId(orderId)
	order.Status = duiBaOrderStatusMap[info.OrderStatus]
	if (order.Status == entity.OrderStatusInTransit || order.Status == entity.OrderStatusComplete) && order.InTransitTime.IsZero() {
		order.InTransitTime = model.NewTime()
	}
	if info.FinishTime.ToInt() > 0 {
		order.CompletedTime = model.Time{Time: time.UnixMilli(info.FinishTime.ToInt())}
	}
	return &order, srv.repo.Save(&order)
}
func (srv OrderService) CreateOrderOfDuiBa(orderId string, info duibaApi.OrderInfo) (*entity.Order, error) {
	order := entity.Order{
		OrderId:          orderId,
		OpenId:           info.Uid,
		TotalCost:        int(info.TotalCredits.ToInt()),
		Status:           duiBaOrderStatusMap[info.OrderStatus],
		PaidTime:         model.Time{Time: time.UnixMilli(info.CreateTime.ToInt())},
		OrderType:        entity.OrderTypePurchase,
		Source:           entity.OrderSourceDuiBa,
		ThirdOrderNo:     info.OrderNum,
		OrderReferenceId: model.NullString(fmt.Sprintf("%d%d", time.Now().Unix(), int(rand.Float64()*10000))),
	}
	if (order.Status == entity.OrderStatusInTransit || order.Status == entity.OrderStatusComplete) && order.InTransitTime.IsZero() {
		order.InTransitTime = model.NewTime()
	}
	if info.FinishTime.ToInt() > 0 {
		order.CompletedTime = model.Time{Time: time.UnixMilli(info.FinishTime.ToInt())}
	}

	orderItemList := make([]entity.OrderItem, 0)
	duibaOrderItemList := info.OrderItemList.OrderItemList()
	for _, duibaOrderItem := range duibaOrderItemList {
		orderItemList = append(orderItemList, entity.OrderItem{
			OrderId: orderId,
			ItemId:  "duiba-" + duibaOrderItem.MerchantCode,
			Count:   int(duibaOrderItem.Quantity.ToInt()),
			Cost:    int(duibaOrderItem.PerCredit.ToInt()),
		})
	}

	return &order, srv.repo.SubmitOrder(&order, &orderItemList)
}
func (srv OrderService) afterCreateOrder(param submitOrderParam, user *entity.User, order *entity.Order, orderItems []entity.OrderItem) {
	participateEvent(param, user, order, orderItems)
	generateBadgeFromOrderItems(param, user, order, orderItems)
}

func participateEvent(param submitOrderParam, user *entity.User, order *entity.Order, orderItems []entity.OrderItem) {
	participateEventParams := make([]event.ParticipateEventParam, 0)
	for _, orderItem := range orderItems {
		participateEventParams = append(participateEventParams, event.ParticipateEventParam{
			ProductItemId: orderItem.ItemId,
			Count:         orderItem.Count,
		})
	}
	err := event.DefaultEventParticipationService.ParticipateEvent(*user, participateEventParams)
	if err != nil {
		app.Logger.Errorf("下订单后参加活动参与失败 %d %+v", user.ID, participateEventParams)
	}
}
func generateBadgeFromOrderItems(param submitOrderParam, user *entity.User, order *entity.Order, orderItems []entity.OrderItem) {
	for _, orderItem := range orderItems {
		cert, err := DefaultCertificateService.FindCertificate(FindCertificateBy{
			ProductItemId: orderItem.ItemId,
		})
		if err != nil {
			app.Logger.Error("发放证书失败 查询证书异常", orderItem.ItemId, err)
			continue
		}

		if cert.ID == 0 {
			continue
		}

		for i := 0; i < orderItem.Count; i++ {
			_, err := DefaultBadgeService.GenerateBadge(GenerateBadgeParam{
				OpenId:        user.OpenId,
				CertificateId: cert.CertificateId,
				ProductItemId: orderItem.ItemId,
				OrderId:       order.OrderId,
				Partnership:   param.PartnershipType,
			})
			if err != nil {
				app.Logger.Error("发放证书失败", orderItem.ItemId, err)
				continue
			}
		}
	}
}
func (srv OrderService) SubmitOrderForEvent(param srv_types.SubmitOrderForEventParam) (*srv_types.SubmitOrderForEventResult, error) {
	ev, err := event.DefaultEventService.FindEvent(event.FindEventParam{
		EventId: param.EventId,
		Active:  sql.NullBool{Bool: true, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	if ev.ID == 0 {
		return nil, errno.ErrRecordNotFound.WithCaller()
	}
	if ev.ProductItemId == "" {
		return nil, errors.New("项目未启用,请稍后再试")
	}
	if !ev.StartTime.IsZero() && ev.StartTime.After(time.Now()) {
		return nil, errno.ErrCommon.WithMessage("项目未开始")
	}

	if !ev.EndTime.IsZero() && ev.EndTime.Before(time.Now()) {
		return nil, errno.ErrCommon.WithMessage("项目已结束")
	}

	order, err := srv.SubmitOrder(SubmitOrderParam{
		Order: SubmitOrder{
			UserId:    param.UserId,
			OrderType: entity.OrderTypePurchase,
		},
		Items: []SubmitOrderItem{
			{
				ItemId: ev.ProductItemId,
				Count:  1,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	badge, err := DefaultBadgeService.FindBadge(srv_types.FindBadgeParam{
		OrderId: order.OrderId,
	})
	if err != nil {
		return nil, err
	}

	code := util2.UUID()

	app.Redis.Set(context.Background(), config.RedisKey.BadgeImageCode+code, badge.ID, time.Minute*5)
	return &srv_types.SubmitOrderForEventResult{
		CertificateNo: badge.Code,
		UploadCode:    code,
	}, nil
}
func (srv OrderService) GetPageFullOrder(dto srv_types.GetPageFullOrderDTO) ([]entity.OrderWithGood, int64, error) {
	orderDO := repotypes.GetPageFullOrderDO{}
	if err := util.MapTo(dto, &orderDO); err != nil {
		return nil, 0, err
	}
	return srv.repo.GetPageFullOrder(orderDO)
}
