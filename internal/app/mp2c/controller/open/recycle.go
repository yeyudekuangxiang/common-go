package open

import (
	context2 "context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/activity/cmd/rpc/activity/activity"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/point/cmd/rpc/point"
	"mio/config"
	"mio/internal/app/mp2c/controller/api"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/queue/producer/recyclepdr"
	"mio/internal/pkg/queue/types/message/recyclemsg"
	"mio/internal/pkg/queue/types/routerkey"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/platform/ccring"
	"mio/internal/pkg/service/platform/recycle"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	"mio/internal/pkg/util/limit"
	"mio/internal/pkg/util/platform"
	"mio/pkg/errno"
	"strconv"
	"strings"
	"time"
)

var DefaultRecycleController = RecycleController{}

type RecycleController struct {
}

func (ctr RecycleController) OolaOrderSync(c *gin.Context) (gin.H, error) {
	// type oola
	form := api.RecyclePushForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	ctx := context.NewMioContext()

	if form.Type != "1" {
		return nil, errno.ErrCommon.WithMessage("非回收订单")
	}

	//查询 渠道信息
	scene := service.DefaultBdSceneService.FindByCh("oola")
	if scene.Key == "" || scene.Key == "e" {
		app.Logger.Info("渠道查询失败", form)
		return nil, errno.ErrCommon.WithMessage("渠道查询失败")
	}

	dst := make(map[string]interface{}, 0)
	err := util.MapTo(&form, &dst)
	if err != nil {
		return nil, err
	}

	//new RecycleService
	RecycleService := recycle.NewRecycleService(ctx)
	TransActionLimitService := service.NewPointTransactionCountLimitService(ctx)
	PointService := service.NewPointService(ctx)
	//校验sign
	if err := RecycleService.CheckSign(dst, scene.Key); err != nil {
		app.Logger.Info("校验sign失败", form)
		return nil, errno.ErrCommon.WithMessage("sign:" + form.Sign + " 验证失败")
	}

	//推送回收订单到消息队列
	defer trackOola(form)

	//通过openid查询用户
	userInfo, exist, err := service.DefaultUserService.GetUser(repository.GetUserBy{OpenId: form.ClientId})
	if err != nil {
		return nil, errno.ErrCommon.WithMessage(err.Error())
	}
	if !exist {
		return nil, errno.ErrUserNotFound
	}

	//查重
	if err = RecycleService.CheckOrder(userInfo.OpenId, "oola"+"#"+form.OrderNo); err != nil {
		fmt.Println("charge 重复提交订单", form)
		app.Logger.Info("charge 重复提交订单", form)
		return nil, errno.ErrCommon.WithMessage("重复提交订单")
	}

	//匹配大类型
	typeName := RecycleService.GetType(form.ProductCategoryName)
	if typeName == "" {
		return nil, errno.ErrCommon.WithMessage("未识别回收分类")
	}

	//入参保存
	defer trackBehaviorInteraction(trackInteractionParam{
		Tp:   string(typeName),
		Data: form,
		Ip:   c.ClientIP(),
	})

	//回调光环
	go ctr.turnPlatform(userInfo, form)

	//本次可得积分
	currPoint, _ := RecycleService.GetPoint(form.Name, form.Qua)
	//本月可获得积分上限
	monthPoint, _ := RecycleService.GetMaxPointByMonth(typeName)

	//最终本次回收可获得积分
	lastPoint, err := TransActionLimitService.CheckMaxPointByMonth(typeName, userInfo.OpenId, currPoint, monthPoint)
	if err != nil {
		return nil, err
	}
	//本次可得减碳量
	currCo2, _ := RecycleService.GetCo2(form.Name, form.Qua)
	carbonString := fmt.Sprintf("%f", currCo2)

	bizId := util.UUID()
	//查询今日该类型获取积分次数
	err = RecycleService.CheckLimit(userInfo.OpenId, form.Name)
	if err == nil {
		//加积分
		_, err = PointService.IncUserPoint(srv_types.IncUserPointDTO{
			OpenId:       userInfo.OpenId,
			Type:         typeName,
			ChangePoint:  lastPoint,
			BizId:        bizId,
			AdditionInfo: form.OrderNo + "#" + strconv.FormatFloat(currCo2, 'f', 2, 64) + "#" + strconv.FormatInt(currPoint, 10) + "#" + form.ClientId,
			Note:         scene.Ch + "#" + form.OrderNo,
		})
		if err != nil {
			app.Logger.Errorf("[oola]旧物回收加积分失败: %s; query:[%v]", err.Error(), form)
		}
	}

	//发碳量
	_, _ = service.NewCarbonTransactionService(context.NewMioContext()).CreateV2(api_types.CreateCarbonTransactionDto{
		OpenId:  userInfo.OpenId,
		UserId:  userInfo.ID,
		Type:    entity.CarbonTransactionType(typeName),
		Value:   currCo2,
		Info:    form.OrderNo + "#" + carbonString + "#" + form.ClientId,
		AdminId: 0,
		Ip:      "",
		BizId:   bizId,
	})
	return gin.H{}, nil
}

func (ctr RecycleController) GetOolaKey(c *gin.Context) (gin.H, error) {
	//查询 渠道信息
	scene := service.DefaultBdSceneService.FindByCh("oola")
	if scene.Key == "" || scene.Key == "e" {
		return nil, errno.ErrCommon.WithMessage("渠道查询失败")
	}
	userInfo := apiutil.GetAuthUser(c)
	oolaPkg := platform.NewOola(context.NewMioContext(), scene.AppId, userInfo.OpenId, scene.Domain, app.Redis)
	oolaPkg.WithHeadImgUrl(userInfo.AvatarUrl)
	oolaPkg.WithUserName(userInfo.Nickname)
	oolaPkg.WithPhone(userInfo.PhoneNumber)
	channelCode, LoginKey, err := oolaPkg.GetToken(scene.Key)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"oolaUserKey": LoginKey,
		"channelCode": channelCode,
	}, nil
}

func (ctr RecycleController) FmyOrderSync(c *gin.Context) (gin.H, error) {
	// type fmy
	form := RecycleFmyForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	ctx := context.NewMioContext()

	if strings.ToUpper(form.Data.Status) != "COMPLETE" {
		return nil, errno.ErrCommon.WithMessage("订单未完成")
	}
	//查询 渠道信息
	scene := service.DefaultBdSceneService.FindByCh("fmy")
	if scene.Key == "" || scene.Key == "e" {
		app.Logger.Info("渠道查询失败", form)
		return nil, errno.ErrCommon.WithMessage("渠道查询失败")
	}
	dst := recycle.FmySignParams{}
	err := util.MapTo(&form, &dst)
	if err != nil {
		return nil, err
	}

	//new RecycleService
	RecycleService := recycle.NewRecycleService(ctx)
	//校验sign
	if err := RecycleService.CheckFmySign(dst, scene.AppId, scene.Key); err != nil {
		app.Logger.Info("校验sign失败", form)
		return nil, errno.ErrCommon.WithMessage("sign:" + form.Sign + " 验证失败")
	}

	defer trackFmy(form)

	if strings.ToUpper(form.Data.Status) != "COMPLETE" {
		return nil, errno.ErrCommon.WithMessage("订单未完成")
	}

	TransActionLimitService := service.NewPointTransactionCountLimitService(ctx)
	PointService := service.NewPointService(ctx)

	//通过phone_number查询用户
	userInfo, exist, err := service.DefaultUserService.GetUser(repository.GetUserBy{
		Source: entity.UserSourceMio,
		Mobile: form.Data.Phone,
	})

	if err != nil {
		return nil, errno.ErrUserNotFound.WithMessage(err.Error())
	}

	if !exist {
		return nil, errno.ErrUserNotFound
	}

	//默认只有衣物
	typeName := entity.POINT_FMY_RECYCLING_CLOTHING
	typeText := RecycleService.GetText(typeName)

	//幂等 检查重复订单
	if err = RecycleService.CheckOrder(userInfo.OpenId, scene.Ch+"#"+form.Data.OrderSn); err != nil {
		return nil, err
	}

	//入参保存
	defer trackBehaviorInteraction(trackInteractionParam{
		Tp:   string(typeName),
		Data: form,
		Ip:   c.ClientIP(),
	})

	//本次可得积分
	currPoint, _ := RecycleService.GetPoint(typeText, form.Data.Weight)
	//本次可得减碳量
	currCo2, _ := RecycleService.GetCo2(typeText, form.Data.Weight)
	carbonString := fmt.Sprintf("%f", currCo2)

	//本月可获得积分上限
	monthPoint, _ := RecycleService.GetMaxPointByMonth(typeName)

	//最终本次回收可获得积分
	lastPoint, err := TransActionLimitService.CheckMaxPointByMonth(typeName, userInfo.OpenId, currPoint, monthPoint)
	if err != nil {
		return nil, err
	}

	//查询今日该类型获取积分次数
	bizId := util.UUID()
	err = RecycleService.CheckLimit(userInfo.OpenId, typeText)
	if err == nil {
		//加积分
		_, err = PointService.IncUserPoint(srv_types.IncUserPointDTO{
			OpenId:       userInfo.OpenId,
			Type:         entity.POINT_FMY_RECYCLING_CLOTHING,
			ChangePoint:  lastPoint,
			BizId:        bizId,
			AdditionInfo: form.Data.OrderSn + "#" + strconv.FormatFloat(currCo2, 'E', -1, 64) + "#" + strconv.FormatInt(currPoint, 10) + "#" + form.Data.Phone,
			Note:         scene.Ch + "#" + form.Data.OrderSn,
		})
		if err != nil {
			app.Logger.Errorf("[fmy]旧物回收加积分失败: %s; query:[%v]", err.Error(), form)
		}
	}

	//发碳量
	_, _ = service.NewCarbonTransactionService(context.NewMioContext()).CreateV2(api_types.CreateCarbonTransactionDto{
		OpenId:  userInfo.OpenId,
		UserId:  userInfo.ID,
		Type:    entity.CarbonTransactionType(typeName),
		Value:   currCo2,
		Info:    form.Data.OrderSn + "#" + carbonString + "#" + form.Data.Phone,
		AdminId: 0,
		Ip:      "",
		BizId:   bizId,
	})

	return gin.H{}, nil
}

//回调的回调
func (ctr RecycleController) turnPlatform(user *entity.User, form api.RecyclePushForm) {
	//绿喵回调ccring
	sceneUser := repository.DefaultBdSceneUserRepository.FindPlatformUser(user.OpenId, "ccring")
	if sceneUser.ID != 0 && sceneUser.PlatformKey == "ccring" {
		ccringScene := service.DefaultBdSceneService.FindByCh("ccring")
		if ccringScene.ID == 0 {
			app.Logger.Info("ccring 渠道查询失败")
			return
		}
		ccRingService := ccring.NewCCRingService("dsaflsdkfjxcmvoxiu123moicuvhoi123", ccringScene.Domain, "/api/cc-ring/external/recycle",
			ccring.WithCCRingOrderNum(form.OrderNo),
			ccring.WithCCRingMemberId(sceneUser.PlatformUserId),
			ccring.WithCCRingProductCategoryName(form.ProductCategoryName),
			ccring.WithCCRingName(form.Name),
			ccring.WithCCRingQua(form.Qua),
		)
		//记录
		one := repository.DefaultBdSceneCallbackRepository.FindOne(repository.GetSceneCallback{
			PlatformKey:    sceneUser.PlatformKey,
			PlatformUserId: sceneUser.PlatformUserId,
			OpenId:         sceneUser.OpenId,
			BizId:          form.OrderNo,
			SourceKey:      "oola",
		})
		if one.ID == 0 {
			body, _ := ccRingService.CallBack()
			err := repository.DefaultBdSceneCallbackRepository.Save(entity.BdSceneCallback{
				PlatformKey:    sceneUser.PlatformKey,
				PlatformUserId: sceneUser.PlatformUserId,
				OpenId:         sceneUser.OpenId,
				BizId:          form.OrderNo,
				SourceKey:      "oola",
				Body:           body,
				CreatedAt:      time.Now(),
			})
			if err != nil {
				return
			}
		}
	}
	return
}

func (ctr RecycleController) Recycle(c *gin.Context) (gin.H, error) {
	form := recycleReq{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	ctx := context.NewMioContext()

	//查询 渠道信息
	scene := service.DefaultBdSceneService.FindByCh(form.Ch)
	if scene.Key == "" || scene.Key == "e" {
		return nil, errno.ErrChannelNotFound.WithMessage("渠道查询失败")
	}

	params := make(map[string]interface{}, 0)
	err := util.MapTo(&form, &params)
	if err != nil {
		return nil, errno.ErrCommon
	}

	//校验sign
	delete(params, "sign")
	if sign := platform.Encrypt(params, scene.Key, "&", "md5"); sign != form.Sign {
		return nil, errno.ErrValidation.WithMessage(fmt.Sprintf("sign:%s 验证失败", form.Sign))
	}

	defer trackRecycle(form)

	//校验用户
	uid, err := strconv.ParseInt(form.MemberId, 10, 64)
	if err != nil {
		return nil, errno.ErrCommon.WithMessage(err.Error())
	}
	userInfo, b, err := service.DefaultUserService.GetUserByID(uid)
	if err != nil {
		return nil, errno.ErrCommon.WithMessage(err.Error())
	}

	if !b {
		return nil, errno.ErrUserNotFound
	}

	//校验重复订单
	RecycleService := recycle.NewRecycleService(ctx)
	if err = RecycleService.CheckOrder(userInfo.OpenId, scene.Ch+"#"+form.OrderNo); err != nil {
		return nil, errno.ErrExisting.WithMessage(fmt.Sprintf("重复订单:%s", form.OrderNo))
	}

	pt := RecycleService.GetPointType(scene.Ch)

	defer trackBehaviorInteraction(trackInteractionParam{
		Tp:   string(pt),
		Data: form,
		Ip:   c.ClientIP(),
	})

	//计算积分
	PointService := service.NewPointService(ctx)
	currPoint, _ := RecycleService.GetPointV2(form.Category, form.Number, form.Name) //本次可得积分
	currCo2, _ := RecycleService.GetCo2V2(form.Category, form.Number, form.Name)     //本次可得减碳量

	//发碳量
	ct := RecycleService.GetCarbonType(scene.Ch)
	_, err = service.NewCarbonTransactionService(context.NewMioContext()).CreateV2(api_types.CreateCarbonTransactionDto{
		OpenId: userInfo.OpenId,
		UserId: userInfo.ID,
		Type:   ct,
		Value:  currCo2,
		Info:   fmt.Sprint(params),
		BizId:  form.OrderNo,
	})
	if err != nil {
		app.Logger.Errorf(fmt.Sprintf("[%s]旧物回收加减碳量失败: %s", form.Ch, err.Error()))
	}

	//每日次数限制
	keyPrefix := fmt.Sprintf("%s:%s:", config.RedisKey.NumberLimit, form.Ch)
	PeriodLimit := limit.NewPeriodLimit(int(time.Hour.Seconds()*24), scene.Override, app.Redis, keyPrefix, limit.PeriodAlign())
	resNumber, err := PeriodLimit.TakeCtx(ctx.Context, form.MemberId)

	if err != nil {
		return nil, errno.ErrInternalServer
	}

	if resNumber != 1 && resNumber != 2 {
		return nil, errno.ErrLimit.WithMessage("超过每日次数上限")
	}

	//每月分数上限 (按分类分别计算)
	maxPoint, err := RecycleService.GetMaxPoint(form.Category)
	if err != nil {
		return nil, errno.ErrCommon.WithMessage(err.Error())
	}

	now := time.Now()
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	lastDay := firstDay.AddDate(0, 1, 0).Add(-time.Nanosecond)
	expired := lastDay.Unix() - time.Now().Unix()
	keyPrefix = fmt.Sprintf("%s:%s:%s:", config.RedisKey.PointMonthLimit, form.Ch, form.Category)

	QuantityLimit := limit.NewQuantityLimit(int(expired), maxPoint, app.Redis, keyPrefix)
	current, err := QuantityLimit.TakeCtx(ctx.Context, form.MemberId, int(currPoint))
	if err != nil {
		return nil, errno.ErrInternalServer
	}
	if current != 0 {
		//加积分
		_, err = PointService.IncUserPoint(srv_types.IncUserPointDTO{
			OpenId:       userInfo.OpenId,
			Type:         pt,
			ChangePoint:  current,
			BizId:        form.OrderNo,
			AdditionInfo: fmt.Sprint(params),
			Note:         scene.Ch + "#" + form.OrderNo,
		})
		if err != nil {
			app.Logger.Errorf(fmt.Sprintf("[%s]旧物回收加积分失败: %s; query:%v", form.Ch, err.Error(), fmt.Sprint(params)))
		}
	}

	//活动
	go func(userId int64, openId, Ch, orderNo, memberId, orderCreateTime, orderCompleteTime string) {
		defer func() {
			if err := recover(); err != nil {
				app.Logger.Errorf("旧物回收活动失败:%v", err)
			}
		}()
		ctr.incPointForActivity(ctx, incPointForActivityParams{
			OpenId:            openId,
			UserId:            userId,
			ActivityCode:      Ch,
			BizId:             orderNo,
			BizName:           memberId,
			OrderCompleteTime: orderCompleteTime,
			OrderCreateTime:   orderCreateTime,
		})
	}(userInfo.ID, userInfo.OpenId, scene.Ch, form.OrderNo, form.MemberId, form.CreateTime, form.CompleteTime)
	return gin.H{}, nil
}
func trackRecycle(req recycleReq) {
	var rk routerkey.RecycleRouterKey
	switch req.Ch {
	case "loverecycle":
		rk = routerkey.RecycleAHS
	case "sshs":
		rk = routerkey.RecycleSSHS
	case "ddyx":
		rk = routerkey.RecycleDDYX
	default:
		return
	}
	_ = recyclepdr.Recycle(rk, recyclemsg.RecycleInfo{
		Ch:           req.Ch,
		OrderNo:      req.OrderNo,
		MemberId:     req.MemberId,
		Name:         req.Name,
		Category:     req.Category,
		Number:       req.Number,
		CreateTime:   req.CreateTime,
		CompleteTime: req.CompleteTime,
		T:            req.T,
		Sign:         req.Sign,
	})
}
func trackFmy(form RecycleFmyForm) {
	_ = recyclepdr.Recycle(routerkey.RecycleFMY, recyclemsg.RecycleFmyInfo{
		AppId:          form.AppId,
		NotificationAt: form.NotificationAt,
		Data: recyclemsg.RecycleFmyData{
			OrderSn:          form.Data.OrderSn,
			Status:           form.Data.Status,
			Weight:           form.Data.Weight,
			Reason:           form.Data.Reason,
			CourierRealName:  form.Data.CourierRealName,
			CourierPhone:     form.Data.CourierPhone,
			CourierJobNumber: form.Data.CourierJobNumber,
			Waybill:          form.Data.Waybill,
			Phone:            form.Data.Phone,
		},
		Sign: form.Sign,
	})
}
func trackOola(form api.RecyclePushForm) {
	_ = recyclepdr.Recycle(routerkey.RecycleOOLA, recyclemsg.RecycleOolaInfo{
		Type:                form.Type,
		OrderNo:             form.OrderNo,
		Name:                form.Name,
		OolaUserId:          form.OolaUserId,
		ClientId:            form.ClientId,
		CreateTime:          form.CreateTime,
		CompletionTime:      form.CompletionTime,
		BeanNum:             form.BeanNum,
		Sign:                form.Sign,
		ProductCategoryName: form.ProductCategoryName,
		Qua:                 form.Qua,
		Unit:                form.Unit,
	})
}
func (ctr RecycleController) incPointForActivity(ctx context2.Context, params incPointForActivityParams) {
	//检查活动
	result, err := app.RpcService.ActivityRpcSrv.ActiveRule(ctx, &activity.ActiveRuleReq{
		Code: params.ActivityCode,
	})
	if err != nil {
		app.Logger.Errorf("用户[%s]参加活动[%s]失败: %s", params.OpenId, params.ActivityCode, err.Error())
		return
	}
	if !result.GetExist() {
		app.Logger.Errorf("用户[%s]参加活动[%s]失败: %s", params.OpenId, params.ActivityCode, "未查询到有效活动规则")
		return
	}
	rule := result.GetActivityRule()
	if rule.GetNumPoint() == 0 {
		return
	}
	//查看用户是否参与过活动并且是否已经领取过积分
	membres, err := app.RpcService.ActivityRpcSrv.Members(ctx, &activity.MembersReq{
		ActivityId: rule.ActivityId,
		UserId:     params.UserId,
	})
	if err != nil {
		app.Logger.Errorf("用户[%s]参加活动[%s]失败: %s", params.OpenId, params.ActivityCode, err.Error())
		return
	}
	if membres.GetMember().GetStatus() == 1 {
		app.Logger.Errorf("用户[%s]参加活动[%s]失败: %s", params.OpenId, params.ActivityCode, "已经参加过活动并且奖励已经领取")
		return
	}
	//三天限制
	parseInt, err := strconv.ParseInt(params.OrderCreateTime, 10, 64)
	if err != nil {
		app.Logger.Errorf("用户[%s]参加活动[%s]失败: %s", params.OpenId, params.ActivityCode, "时间解析失败")
	}
	if time.UnixMilli(parseInt).Sub(time.UnixMilli(membres.GetMember().GetCreatedAt())).Hours() > 72.0 {
		app.Logger.Errorf("用户[%s]参加活动[%s]失败: %s", params.OpenId, params.ActivityCode, "超出3天时间限制")
		return
	}
	//记录次数
	expired := (rule.GetEndTime() - rule.GetStartTime()) / 1000 //秒级
	QuantityLimit := limit.NewQuantityLimit(int(expired), int(rule.GetNumLimit()), app.Redis, config.RedisKey.NumberLimit)
	current, err := QuantityLimit.TakeCtx(ctx, rule.GetActivityCode(), 1)
	if err != nil {
		app.Logger.Errorf("用户[%s]参加活动[%s]-[%s], 次数记录失败: %s", params.OpenId, params.ActivityCode, rule.GetTitle(), err.Error())
	}
	if current == 0 {
		//达到上限 不再赠送积分
		app.Logger.Errorf("用户[%s]参加活动[%s]-[%s], 活动奖励发放次数达到上限", params.OpenId, params.ActivityCode, rule.GetTitle())
		return
	}
	//发放积分
	_, err = app.RpcService.PointRpcSrv.IncPoint(ctx, &point.IncPointReq{
		Openid:      params.OpenId,
		Type:        string(entity.POINT_PLATFORM),
		BizId:       params.BizId,
		BizName:     params.BizName,
		ChangePoint: uint64(rule.NumPoint),
		Node:        "活动赠送积分",
	})
	if err != nil {
		app.Logger.Errorf("用户[%s]参加活动[%s]-[%s], 积分发放失败: %s", params.OpenId, params.ActivityCode, rule.GetTitle(), err.Error())
		//加积分失败次数-1
		app.Redis.DecrBy(ctx, config.RedisKey.NumberLimit+rule.GetActivityCode(), 1)
		return
	}
	//更新用户参与记录
	_, err = app.RpcService.ActivityRpcSrv.MemberUpdate(ctx, &activity.MemberUpdateReq{
		ActivityId: rule.GetActivityId(),
		UserId:     params.UserId,
		Status:     1,
	})
	if err != nil {
		app.Logger.Errorf("用户[%s]参加活动[%s]-[%s], 更新用户活动状态失败: %s", params.OpenId, params.ActivityCode, rule.GetTitle(), err.Error())
		return
	}
	//更新发放次数
	_, err = app.RpcService.ActivityRpcSrv.IncNumSended(ctx, &activity.IncNumSendedReq{Id: rule.GetId()})
	if err != nil {
		app.Logger.Errorf("用户[%s]参加活动[%s]-[%s], 更新发放次数失败: %s", params.OpenId, params.ActivityCode, rule.GetTitle(), err.Error())
		return
	}
	return
}
