package service

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"log"
	"math/rand"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	eevent "mio/internal/pkg/model/entity/event"
	repository2 "mio/internal/pkg/repository"
	repositoryActivity "mio/internal/pkg/repository/activity"
	"mio/internal/pkg/repository/repotypes"
	"mio/internal/pkg/service/event"
	"mio/internal/pkg/service/platform"
	"mio/internal/pkg/service/product"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/service/track"
	util2 "mio/internal/pkg/util"
	duibaApi "mio/pkg/duiba/api/model"
	"mio/pkg/duiba/util"
	"mio/pkg/errno"
	"strconv"
	"time"
)

var DefaultOrderService = NewOrderService(repository2.DefaultOrderRepository, repository2.DefaultOrderRepository, repository2.DefaultUserRepository, repositoryActivity.DefaultGDDonationBookRepository)

func NewOrderService(repo repository2.OrderRepository, repository repository2.OrderRepository, userRepository repository2.UserRepository, repositoryActivity repositoryActivity.GDDonationBookRepository) OrderService {
	return OrderService{repo: repo, repoOrder: repository, repoUser: userRepository, repoGDBook: repositoryActivity}
}

type OrderService struct {
	repo       repository2.OrderRepository
	repoOrder  repository2.OrderRepository
	repoUser   repository2.UserRepository
	repoGDBook repositoryActivity.GDDonationBookRepository
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
	if !util2.DefaultLock.Lock("submitOrder_"+strconv.FormatInt(param.Order.UserId, 5), time.Second*5) {
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

	pointService := NewPointService(context.NewMioContext())
	//检查积分
	point, err := pointService.FindByUserId(param.Order.UserId)
	if err != nil {
		return nil, err
	}
	if point.Balance < int64(param.Order.TotalCost) {
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
	_, err = pointService.DecUserPoint(srv_types.DecUserPointDTO{
		OpenId:       user.OpenId,
		ChangePoint:  int64(param.Order.TotalCost),
		BizId:        util2.UUID(),
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

			_, err = pointService.IncUserPoint(srv_types.IncUserPointDTO{
				OpenId:       user.OpenId,
				ChangePoint:  int64(param.Order.TotalCost),
				BizId:        util2.UUID(),
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
		OrderReferenceId: model.NullString(fmt.Sprintf("%d%d", time.Now().UnixMilli(), int(rand.Float64()*10000))),
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
		OrderReferenceId: model.NullString(fmt.Sprintf("%d%d", time.Now().UnixMilli(), int(rand.Float64()*10000))),
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
			ItemId:  duibaOrderItem.MerchantCode,
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

			//根据商品id，查询证书id
			eventInfo, errEvent := event.DefaultEventService.FindEventAndCate(event.FindEventParam{ProductItemId: orderItem.ItemId})
			var title, cateTitle string
			if errEvent != nil {
				title = ""
				cateTitle = ""
			} else {
				title = eventInfo.Title
				cateTitle = eventInfo.CateTitle
			}
			//证书诸葛打点
			zhuGeAttr := make(map[string]interface{}, 0)
			zhuGeAttr["分类名称"] = cateTitle
			zhuGeAttr["证书id"] = cert.CertificateId
			zhuGeAttr["商品id"] = orderItem.ItemId
			zhuGeAttr["项目名称"] = title

			if err != nil {
				app.Logger.Error("发放证书失败", orderItem.ItemId, err)
				zhuGeAttr["是否失败"] = "操作失败"
				zhuGeAttr["失败原因"] = err.Error()
			} else {
				zhuGeAttr["是否失败"] = "操作成功"
			}
			track.DefaultZhuGeService().Track(config.ZhuGeEventName.UserCertificateSendSuc, user.OpenId, zhuGeAttr)
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

	if err := srv.checkEventLimit(param.UserId, ev.EventId, ev.Limit); err != nil {
		return nil, err
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

	app.Redis.Set(context.NewMioContext(), config.RedisKey.BadgeImageCode+code, badge.ID, time.Minute*5)

	srv.addEventLimit(param.UserId, ev.EventId, ev.Limit)

	srv.afterSubmitEventOrder(param.UserId, ev.EventId, ev.Limit, badge)
	return &srv_types.SubmitOrderForEventResult{
		CertificateNo: badge.Code,
		UploadCode:    code,
	}, nil
}

func (srv OrderService) SubmitOrderForEventGD(param srv_types.SubmitOrderForEventGDParam) (*srv_types.SubmitOrderForEventResult, error) {
	if !util2.DefaultLock.Lock("SubmitOrderForEventGD"+param.OpenId, time.Second*10) {
		return nil, errno.ErrLimit
	}
	defer util2.DefaultLock.UnLock("SubmitOrderForEventGD" + param.OpenId)

	wechatServiceOpenId := param.WechatServiceOpenId
	info := srv.repoUser.GetUserBy(repository2.GetUserBy{OpenId: wechatServiceOpenId})
	wechatServiceUid := info.ID
	if wechatServiceUid == 0 {
		return nil, errors.New("不满足领取条件")
	}
	wechatServiceUser := srv.repoGDBook.GetUserBy(repositoryActivity.FindRecordBy{UserId: wechatServiceUid})
	if wechatServiceUser.UserId == 0 {
		return nil, errors.New("您不满足领取条件哦")
	}

	openid := param.OpenId
	//判断是否领取过证书
	var ItemIdSlice = []string{"cbddf0af60f402f717b0987b79709209", "b00064a760f400a42850b68e1f783c22"}
	orderTotal := srv.repoOrder.GetOrderTotalByItemId(repotypes.GetOrderTotalByItemIdDO{
		Openid:      openid,
		ItemIdSlice: ItemIdSlice})
	if orderTotal >= 1 {
		return nil, errors.New("您已经领取过证书了")
	}
	order, errorOrder := srv.SubmitOrderForEvent(srv_types.SubmitOrderForEventParam{UserId: param.UserId, EventId: param.EventId})
	if errorOrder != nil {
		return nil, errorOrder
	}

	//发放积分
	point := 500
	pointService := NewPointService(context.NewMioContext())
	_, errInc := pointService.IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       openid,
		Type:         entity.POINT_PARTNERSHIP,
		ChangePoint:  int64(point),
		BizId:        util2.UUID(),
		AdditionInfo: order.UploadCode + "#" + order.UploadCode + "#" + strconv.Itoa(point),
		Note:         order.UploadCode + "#" + order.UploadCode,
	})
	if errInc != nil {
		fmt.Println("广东教育学会，发证书 加积分失败，失败原因", errInc.Error())
		return nil, errors.New(errInc.Error())
	}
	return order, nil
}

func (srv OrderService) afterSubmitEventOrder(userId int64, evId string, limit eevent.EventLimit, badge *entity.Badge) {
	app.Logger.Info("提交兑换订单后", evId, config.Constants.StarCouponEventId)
	if config.Constants.StarCouponEventId == evId {
		log.Println("发送星星充电券", evId)
		srv.sendEventStarCoupon(userId, evId)
	}

	return
}
func (srv OrderService) sendEventStarCoupon(userId int64, evId string) {
	userInfo, exist, err := DefaultUserService.GetUserByID(userId)
	if err != nil {
		app.Logger.Error("兑换证书发星星券失败,查询用户信息异常", err, userId, evId)
		return
	}
	if !exist {
		app.Logger.Error("兑换证书发星星券失败,未查询到用户信息", userId, evId)
		return
	}
	if userInfo.PhoneNumber == "" {
		app.Logger.Error("兑换证书发星星券失败,用户未绑定手机号", userId, evId)
		return
	}

	starChargeService := platform.NewStarChargeService(context.NewMioContext())
	token, err := starChargeService.GetAccessToken()
	if err != nil {
		app.Logger.Error("兑换证书发星星券失败,获取星星token失败", err, userId, evId)
		return
	}
	err = starChargeService.SendCoupon(userInfo.OpenId, userInfo.PhoneNumber, starChargeService.ProvideId, token)
	if err != nil {
		app.Logger.Error("兑换证书发星星券失败", err, userId, evId)
		return
	}
	app.Logger.Info("发送星星充电券成功", userInfo.OpenId, userInfo.PhoneNumber, evId)
}
func (srv OrderService) addEventLimit(userId int64, evId string, limit eevent.EventLimit) {
	redisKey := fmt.Sprintf("%s:%d%s", config.RedisKey.EventLimit, userId, evId)

	d, c, err := limit.Parse()
	if err != nil {
		app.Logger.Error(limit, err)
		return
	}

	if d == 0 {
		return
	}

	usedCount, err := app.Redis.Get(context.NewMioContext(), redisKey).Int64()
	if err != nil && err != redis.Nil {
		app.Logger.Error(limit, err)
	}

	if usedCount >= c {
		app.Logger.Error("超过兑换次数上限", userId, limit, usedCount)
	}

	if err == redis.Nil {
		app.Redis.SetNX(context.NewMioContext(), redisKey, 1, d)
	} else {
		app.Redis.Incr(context.NewMioContext(), redisKey)
	}

	return
}
func (srv OrderService) checkEventLimit(userId int64, evId string, limit eevent.EventLimit) error {
	redisKey := fmt.Sprintf("%s:%d%s", config.RedisKey.EventLimit, userId, evId)

	d, c, err := limit.Parse()
	if err != nil {
		app.Logger.Error(limit, err)
		return errno.ErrInternalServer
	}
	if d == 0 {
		return nil
	}

	usedCount, err := app.Redis.Get(context.NewMioContext(), redisKey).Int64()
	if err != nil && err != redis.Nil {
		app.Logger.Error(limit, err)
	}
	if usedCount >= c {
		return errno.ErrCommon.WithMessage("已达到兑换次数上限")
	}
	return nil
}
func (srv OrderService) GetPageFullOrder(dto srv_types.GetPageFullOrderDTO) ([]entity.OrderWithGood, int64, error) {
	orderDO := repotypes.GetPageFullOrderDO{}
	if err := util.MapTo(dto, &orderDO); err != nil {
		return nil, 0, err
	}
	return srv.repo.GetPageFullOrder(orderDO)
}
