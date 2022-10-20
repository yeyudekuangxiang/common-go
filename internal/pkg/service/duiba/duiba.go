package duiba

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/platform/jhx"
	"mio/internal/pkg/service/product"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/encrypt"
	"mio/pkg/duiba"
	duibaApi "mio/pkg/duiba/api/model"
	"mio/pkg/errno"
	"time"
)

var DefaultDuiBaService DuiBaService

func NewDuiBaService(client *duiba.Client) DuiBaService {
	return DuiBaService{
		client: client,
	}
}

// InitDefaultDuibaService 配置文件加载后调用此方法初始化默认兑吧服务
func InitDefaultDuibaService() {
	client := duiba.NewClient(config.Config.DuiBa.AppKey, config.Config.DuiBa.AppSecret)
	DefaultDuiBaService = NewDuiBaService(client)
}

type DuiBaService struct {
	client *duiba.Client
}

// AutoLogin 获取免登陆地址
func (srv DuiBaService) AutoLogin(param service.AutoLoginParam) (string, error) {
	userInfo, err := service.DefaultUserService.GetUserById(param.UserId)
	if err != nil {
		return "", err
	}
	b, err := service.NewPointService(context.NewMioContext()).FindByUserId(param.UserId)
	if err != nil {
		return "", err
	}
	return srv.client.AutoLogin(duiba.AutoLoginParam{
		Uid:      userInfo.OpenId,
		Credits:  int64(b.Balance),
		Redirect: param.Path,
		DCustom:  param.DCustom,
		Transfer: param.Transfer,
		SignKeys: param.SignKeys,
	})
}
func (srv DuiBaService) AutoLoginOpenId(param service.AutoLoginOpenIdParam) (string, error) {
	/*userInfo, err := DefaultUserService.GetUserById(param.UserId)
	if err != nil {
		return "", err
	}*/
	b, err := service.NewPointService(context.NewMioContext()).FindByUserId(param.UserId)
	if err != nil {
		return "", err
	}
	return srv.client.AutoLogin(duiba.AutoLoginParam{
		Uid:      param.OpenId,
		Credits:  int64(b.Balance),
		Redirect: param.Path,
		DCustom:  param.DCustom,
		Transfer: param.Transfer,
		SignKeys: param.SignKeys,
		Vip:      param.Vip,
	})
}

var duibaTypeToPointType = map[duibaApi.ExchangeType]entity.PointTransactionType{
	duibaApi.ExchangeTypeAlipay:    entity.POINT_DUIBA_ALIPAY,
	duibaApi.ExchangeTypeQB:        entity.POINT_DUIBA_QB,
	duibaApi.ExchangeTypeCoupon:    entity.POINT_DUIBA_COUPON,
	duibaApi.ExchangeTypeObject:    entity.POINT_DUIBA_OBJECT,
	duibaApi.ExchangeTypePhoneBill: entity.POINT_DUIBA_PHONEBILL,
	duibaApi.ExchangeTypePhoneFlow: entity.POINT_DUIBA_PHONEFLOW,
	duibaApi.ExchangeTypeVirtual:   entity.POINT_DUIBA_VIRTUAL,
	duibaApi.ExchangeTypeGame:      entity.POINT_DUIBA_GAME,
	duibaApi.ExchangeTypeHdTool:    entity.POINT_DUIBA_HDTOOL,
	duibaApi.ExchangeTypeHdSign:    entity.POINT_DUIBA_SIGN,
}

// ExchangeCallback 扣积分回调
func (srv DuiBaService) ExchangeCallback(form duibaApi.Exchange) (*service.ExchangeCallbackResult, error) {
	userInfo, err := service.DefaultUserService.GetUserBy(repository.GetUserBy{
		OpenId: form.Uid,
	})
	if err != nil {
		return nil, err
	}
	if userInfo.ID == 0 {
		return nil, errno.ErrCommon.WithMessage("用户信息不存在")
	}
	data, err := json.Marshal(form)
	if err != nil {
		app.Logger.Errorf("%+v %v", form, err)
		return nil, errors.New("系统异常,请联系管理员")
	}

	pointType := duibaTypeToPointType[form.Type]

	bizId := util.UUID()
	balance, err := service.NewPointService(context.NewMioContext()).DecUserPoint(srv_types.DecUserPointDTO{
		OpenId:       form.Uid,
		Type:         pointType,
		BizId:        bizId,
		ChangePoint:  form.Credits,
		AdditionInfo: string(data),
	})
	if err != nil {
		app.Logger.Errorf("%+v %v", form, err)
		return nil, errors.New("系统异常,请联系管理员")
	}

	return &service.ExchangeCallbackResult{
		BizId:   bizId,
		Credits: balance,
	}, nil
}

// ExchangeResultNoticeCallback 积分兑换结果回调
func (srv DuiBaService) ExchangeResultNoticeCallback(form duibaApi.ExchangeResult) error {
	if form.Success {
		return nil
	}

	if form.BizId == "" {
		app.Logger.Error("订单异常", form)
		return nil
	}

	userInfo, err := service.DefaultUserService.GetUserBy(repository.GetUserBy{
		OpenId: form.Uid,
	})
	if err != nil {
		return err
	}
	if userInfo.ID == 0 {
		return errno.ErrCommon.WithMessage("用户信息不存在")
	}

	pointTranService := service.NewPointTransactionService(context.NewMioContext())
	pt, err := pointTranService.FindBy(repository.FindPointTransactionBy{
		TransactionId: form.BizId,
	})
	if err != nil {
		return err
	}
	if pt.ID == 0 {
		return nil
	}

	data, err := json.Marshal(form)
	if err != nil {
		return err
	}
	refundPoint := -pt.Value

	pointService := service.NewPointService(context.NewMioContext())
	_, err = pointService.IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       form.Uid,
		Type:         entity.POINT_DUIBA_REFUND,
		BizId:        util.UUID(),
		ChangePoint:  refundPoint,
		AdditionInfo: string(data),
	})
	return err
}

func (srv DuiBaService) OrderCallback(form duibaApi.OrderInfo) error {
	user, err := service.DefaultUserService.GetUserByOpenId(form.Uid)
	if err != nil {
		return err
	}
	if user.ID == 0 {
		return errno.ErrUserNotFound.WithCaller()
	}

	orderId := form.DevelopBizId
	if orderId == "" {
		orderId = "duiba-" + form.OrderNum
	}

	orderItemList := form.OrderItemList.OrderItemList()
	for i, orderItem := range orderItemList {
		itemId := ""
		if orderItem.MerchantCode != "" {
			itemId = "duiba-" + orderItem.MerchantCode
		} else {
			itemId = "duiba-" + form.OrderNum + "-" + encrypt.Md5(orderItem.Title)
		}
		orderItemList[i].MerchantCode = itemId

		_, err := product.DefaultProductItemService.CreateOrUpdateProductItem(product.CreateOrUpdateProductItemParam{
			ItemId:   itemId,
			Virtual:  false,
			Title:    orderItem.Title,
			Cost:     int(orderItem.PerCredit.ToInt()),
			ImageUrl: orderItem.SmallImage,
		})
		if err != nil {
			return err
		}
	}
	orderItemData, err := json.Marshal(orderItemList)
	if err != nil {
		return err
	}
	form.OrderItemList = duibaApi.OrderItemListStr(orderItemData)

	_, err = service.DefaultDuiBaOrderService.CreateOrUpdate(orderId, form)
	if err != nil {
		return err
	}
	_, err = service.DefaultOrderService.CreateOrUpdateOrderOfDuiBa(orderId, form)
	return err
}
func (srv DuiBaService) CheckSign(param duiba.Param) error {
	return srv.client.CheckSign(param)
}

var duibaPointAddTypeMap = map[duibaApi.PointAddType]entity.PointTransactionType{
	duibaApi.PointAddTypeGame:       entity.POINT_DUIBA_GAME,
	duibaApi.PointAddTypeSign:       entity.POINT_DUIBA_SIGN,
	duibaApi.PointAddTypeTask:       entity.POINT_DUIBA_TASK,
	duibaApi.PointAddTypeReSign:     entity.POINT_DUIBA_SIGN,       //补签
	duibaApi.PointAddTypePostSale:   entity.POINT_DUIBA_POSTSALE,   //售后退积分
	duibaApi.PointAddTypeCancelShip: entity.POINT_DUIBA_CANCELSHIP, //取消发货
	duibaApi.PointAddTypeHdTool:     entity.POINT_DUIBA_HDTOOL,
}

func (srv DuiBaService) PointAddCallback(form duibaApi.PointAdd) (tranId string, err error) {
	log, err := service.DefaultDuiBaPointAddLogService.FindBy(service.FindDuiBaPointAddLogBy{
		OrderNum: form.OrderNum,
	})
	if err != nil {
		return
	}

	if log.TransactionId != "" {
		return log.TransactionId, nil
	}
	newLog, err := service.DefaultDuiBaPointAddLogService.CreateLog(service.CreateDuiBaPointAddLog{
		Uid:         form.Uid,
		Credits:     form.Credits.ToInt(),
		Type:        form.Type,
		OrderNum:    form.OrderNum,
		SubOrderNum: form.SubOrderNum,
		Timestamp:   form.Timestamp.ToInt(),
		Description: form.Description,
		Ip:          form.IP,
		Sign:        form.Sign,
		AppKey:      form.AppKey,
	})
	if err != nil {
		return
	}

	pointType := duibaPointAddTypeMap[form.Type]
	if pointType == "" {
		pointType = entity.POINT_ADJUSTMENT
	}
	bizId := util.UUID()
	_, err = service.NewPointService(context.NewMioContext()).IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       form.Uid,
		Type:         pointType,
		ChangePoint:  form.Credits.ToInt(),
		BizId:        bizId,
		AdminId:      0,
		AdditionInfo: fmt.Sprintf("log %d", newLog.ID),
	})
	if err != nil {
		return
	}

	err = service.DefaultDuiBaPointAddLogService.UpdateLogTransaction(newLog.ID, bizId)
	if err != nil {
		app.Logger.Errorf("更新DuiBaPointAddLog失败 %d %s", newLog.ID, bizId)
	}

	return bizId, nil
}

var virtualGoodMap = map[string]int{
	"1185f098-ffae-4c4f-ae17-04c739d7664d": 50,
	"87ac16fa-f2c9-4b3a-9008-992353f0ed39": 100,
	"8791a99f-9a66-44a1-b3de-cb4e492b4241": 150,
	"fc0b48e2-056f-40ee-b729-ef76bc3996bf": 200,
	"b8dbe8af-b0cb-4c80-bda1-fdbf127297d4": 500,
	"923e3f57-2cd5-4d3a-9dbd-e79f10f5c65c": 1000,
	"bd9463d6-f81e-44c4-b085-5ded53d8fe34": 1500,
	"6e4a4a84-92af-43d4-b73b-8da09d9647b0": 2000,
	"f233927c-72cc-4df3-8bc3-b4fe6d8a31b9": 2500,
	"7380d8bb-9969-4028-90bc-06a357b40abb": 3000,
	"0e72028b-7c93-4a64-ba0f-2144b08ffef4": 5000,
	"6c2a99fc-4b6f-49da-9b23-b47c11b8494f": 10000,
	"136c0116-df2c-4610-b3c2-aabbec90b75a": 15000,
	"90de1c36-60d7-4e75-b738-79ce2ab9a405": 20000,
}

func (srv DuiBaService) VirtualGoodCallback(form duibaApi.VirtualGood) (orderId string, credit int64, err error) {
	lockKey := fmt.Sprintf("VirtualGoodCallback%s", encrypt.Md5(form.OrderNum+form.Params))
	if !util.DefaultLock.Lock(lockKey, time.Second*10) {
		return "", 0, errno.ErrCommon.WithMessage("操作频繁,请稍后再试")
	}
	defer util.DefaultLock.UnLock(lockKey)

	//幂等
	log, err := DefaultVirtualGoodLogService.FindVirtualGoodLog(FindVirtualGoodLogParam{
		OrderNum: form.OrderNum,
		Params:   form.Params,
	})
	if err != nil {
		return "", 0, err
	}
	pointService := service.NewPointService(context.NewMioContext())
	if log.ID != 0 {
		userPoint, err := pointService.FindByOpenId(form.Uid)
		if err != nil {
			return "", 0, err
		}
		return log.SupplierBizId, userPoint.Balance, nil
	}

	log, err = DefaultVirtualGoodLogService.CreateVirtualGoodLog(form)
	if err != nil {
		return "", 0, err
	}

	switch form.Params {
	case virtualCouponJhx2Yuan:
		err := srv.SendVirtualCoupon(form.OrderNum, form.Uid, form.Params)
		if err != nil {
			app.Logger.Error("发放兑吧虚拟商品优惠券失败", err)
			return "", 0, err
		}
		userPoint, err := pointService.FindByOpenId(form.Uid)
		if err != nil {
			return "", 0, err
		}
		return log.SupplierBizId, userPoint.Balance, nil
	}

	if _, ok := virtualGoodMap[form.Params]; ok {
		point, errSendPoint := srv.SendVirtualGoodPoint(form.OrderNum, form.Uid, form.Params)
		if errSendPoint != nil {
			app.Logger.Error("发放兑吧虚拟商品积分失败", errSendPoint)
			return "", 0, errSendPoint
		} else {
			return log.SupplierBizId, point, nil
		}
	}

	return "", 0, errno.ErrCommon.WithMessage("虚拟商品不存在")
}
func (srv DuiBaService) SendVirtualGoodPoint(orderNum, openid string, productItemId string) (int64, error) {
	point := virtualGoodMap[productItemId]
	if point == 0 {
		return 0, errno.ErrCommon.WithMessage("虚拟商品不存在")
	}
	pointService := service.NewPointService(context.NewMioContext())
	newPoint, err := pointService.IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       openid,
		Type:         entity.POINT_DUIBA_INTEGRAL_RECHARGE,
		ChangePoint:  int64(point),
		BizId:        util.UUID(),
		AdditionInfo: fmt.Sprintf("兑吧虚拟商品兑换 orderNum:%s productItemId:%s", orderNum, productItemId),
	})
	if err != nil {
		return 0, err
	}
	return newPoint, nil
}

const (
	virtualCouponJhx2Yuan = "3323df0ce743a3e55a38c62dbc92eac4"
)

func (srv DuiBaService) SendVirtualCoupon(orderNum, openid, productItemId string) error {
	switch productItemId {
	case virtualCouponJhx2Yuan:
		jhxService := jhx.NewJhxService(context.NewMioContext())
		user, err := service.DefaultUserService.GetUserByOpenId(openid)
		if err != nil {
			return err
		}
		if user.ID == 0 {
			return errno.ErrUserNotFound.WithCaller()
		}
		tradeNo, err := jhxService.SendCoupon(1000, *user)
		println(tradeNo)
		if err != nil {
			return err
		}
		return nil
	}
	app.Logger.Error("未知的虚拟商品类型", orderNum, openid, productItemId)
	return errno.ErrCommon.WithMessage("未知的虚拟商品类型")
}
