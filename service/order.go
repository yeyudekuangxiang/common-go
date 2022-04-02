package service

import (
	"errors"
	"fmt"
	"math/rand"
	"mio/core/app"
	"mio/internal/util"
	"mio/model"
	"mio/model/entity"
	"mio/repository"
	"strconv"
	"time"
)

var DefaultOrderService = NewOrderService(repository.DefaultOrderRepository)

func NewOrderService(repo repository.OrderRepository) OrderService {
	return OrderService{repo: repo}
}

type OrderService struct {
	repo repository.OrderRepository
}

// CalculateAndCheck 计算商品价格并且检查库存
func (srv OrderService) CalculateAndCheck(items []repository.CheckStockItem) (*CalculateProductResult, error) {
	itemIds := make([]string, 0)
	itemMap := make(map[string]repository.CheckStockItem)
	for _, item := range items {
		itemIds = append(itemIds, item.ItemId)
		itemMap[item.ItemId] = item
	}

	productItems := DefaultProductItemService.GetListBy(repository.GetProductItemListBy{
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
	checkItems := make([]repository.CheckStockItem, 0)
	for _, item := range param.Items {
		checkItems = append(checkItems, repository.CheckStockItem{
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
		Items: calculateResult.ItemList,
	})
}

//submitOrder (此方法可自定义需要支付的金额 需谨慎使用 用户下单请使用 OrderService.SubmitOrder 方法创建订单)
func (srv OrderService) submitOrder(param submitOrderParam) (*entity.Order, error) {
	//防止并发
	if !util.Lock.Lock("submitOrder_"+strconv.FormatInt(param.Order.UserId, 10), time.Second*5) {
		return nil, errors.New("操作频率过快,请稍后再试")
	}

	orderSuccess := false
	orderId := util.UUID()

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
	checkStockItems := make([]repository.CheckStockItem, 0)
	for _, item := range param.Items {
		checkStockItems = append(checkStockItems, repository.CheckStockItem{
			ItemId: item.ItemId,
			Count:  item.Count,
		})
	}
	err = DefaultProductItemService.CheckAndLockStock(checkStockItems)
	if err != nil {
		return nil, err
	}
	//下单失败释放库存
	defer func() {
		if orderSuccess == false {
			//释放库存
			err := DefaultProductItemService.UnLockStock(checkStockItems)
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
	order, err := srv.create(orderId, param)
	if err != nil {
		return nil, err
	}

	orderSuccess = true
	return order, nil
}

//直接创建订单(请勿使用此方法创建订单 请使用 OrderService.SubmitOrder 方法创建订单)
func (srv OrderService) create(orderId string, param submitOrderParam) (*entity.Order, error) {
	var addressId *string
	if param.Order.AddressId != "" {
		addressId = &param.Order.AddressId
	}
	user, err := DefaultUserService.GetUserById(param.Order.UserId)
	if err != nil {
		return nil, err
	}
	if user.ID == 0 || user.OpenId == "" {
		return nil, errors.New("未查找到用户信息,请联系管理员")
	}

	order := &entity.Order{
		OrderId:          orderId,
		AddressId:        addressId,
		OpenId:           user.OpenId,
		TotalCost:        param.Order.TotalCost,
		Status:           entity.OrderStatusPaid,
		PaidTime:         model.NewTime(),
		OrderReferenceId: fmt.Sprintf("%d%d", time.Now().Unix(), int(rand.Float64()*10000)),
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

	return order, srv.repo.SubmitOrder(order, &orderItems)
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
