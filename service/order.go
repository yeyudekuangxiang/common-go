package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
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
func NewOrderServiceByDB(db *gorm.DB) OrderService {
	return NewOrderService(repository.NewOrderRepository(db))
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
func (srv OrderService) submitOrder(param submitOrderParam) (*entity.Order, error) {
	//防止并发
	if !util.Lock.Lock("submitOrder_"+strconv.FormatInt(param.Order.UserId, 10), time.Second*5) {
		return nil, errors.New("操作频率过快,请稍后再试")
	}
	var order *entity.Order

	user, err := DefaultUserService.GetUserById(param.Order.UserId)
	if err != nil {
		return nil, err
	}
	if user.ID == 0 || user.OpenId == "" {
		return nil, errors.New("未查找到用户信息,请联系管理员")
	}

	err = app.DB.Transaction(func(tx *gorm.DB) error {
		orderService := NewOrderServiceByDB(tx)
		orderItemService := NewOrderItemServiceByDB(tx)
		productItemService := NewProductItemServiceByDB(tx)
		orderSuccess := false

		//检查积分
		point, err := DefaultPointService.FindByUserId(param.Order.UserId)
		if err != nil {
			return err
		}
		if point.Balance < param.Order.TotalCost {
			return errors.New("积分不足,无法兑换")
		}

		//检查库存
		checkStockItems := make([]repository.CheckStockItem, 0)
		for _, item := range param.Items {
			checkStockItems = append(checkStockItems, repository.CheckStockItem{
				ItemId: item.ItemId,
				Count:  item.Count,
			})
		}
		err = productItemService.CheckAndLockStock(checkStockItems)
		if err != nil {
			return err
		}

		orderId := util.UUID()

		defer func() {
			if orderSuccess == false {
				//释放库存
				err := productItemService.UnLockStock(checkStockItems)
				if err != nil {
					app.Logger.Errorf("释放库存失败 %+v %+v", checkStockItems, err)
				}

				//返还积分
				_, err = DefaultPointTransactionService.Create(CreatePointTransactionParam{
					OpenId:       user.OpenId,
					Value:        param.Order.TotalCost,
					Type:         entity.POINT_ADJUSTMENT,
					AdditionInfo: `{"orderId":"` + orderId + `"}`,
				})

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
			return err
		}
		//创建订单
		order, err = orderService.create(param.Order)
		if err != nil {
			return err
		}

		//创建订单item
		err = orderItemService.CreateOrderItems(orderId, param.Items)
		if err != nil {
			return err
		}
		orderSuccess = true
		return nil
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (srv OrderService) create(param submitOrder) (*entity.Order, error) {
	var addressId *string
	if param.AddressId != "" {
		addressId = &param.AddressId
	}

	user, err := DefaultUserService.GetUserById(param.UserId)
	if err != nil {
		return nil, err
	}
	if user.ID == 0 || user.OpenId == "" {
		return nil, errors.New("未查找到用户信息,请联系管理员")
	}

	order := &entity.Order{
		OrderId:          util.UUID(),
		AddressId:        addressId,
		OpenId:           user.OpenId,
		TotalCost:        param.TotalCost,
		Status:           entity.OrderStatusPaid,
		PaidTime:         model.NewTime(),
		OrderReferenceId: fmt.Sprintf("%d%d", time.Now().Unix(), int(rand.Float64()*10000)),
		OrderType:        param.OrderType,
	}
	return order, srv.repo.Save(order)
}

// SubmitOrderForGreenMonday 用于greenmonday活动用户下单
func (srv OrderService) SubmitOrderForGreenMonday(param SubmitOrderForGreenParam) (*entity.Order, error) {
	return srv.SubmitOrder(SubmitOrderParam{
		Order: SubmitOrder{
			UserId:    param.UserId,
			AddressId: param.AddressId,
			OrderType: entity.OrderTypePurchase,
		},
		Items: []SubmitOrderItem{
			{
				ItemId: param.ItemId,
				Count:  1,
			},
		},
	})
}
