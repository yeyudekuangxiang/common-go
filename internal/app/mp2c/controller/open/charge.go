package open

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/activity/cmd/rpc/activity/activity"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/point/cmd/rpc/point"
	"mio/config"
	"mio/internal/app/mp2c/controller/api"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/app"
	mioContext "mio/internal/pkg/core/context"

	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/platform/ccring"
	"mio/internal/pkg/service/platform/star_charge"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	"mio/internal/pkg/util/limit"
	"mio/pkg/errno"
	"strconv"
	"time"
)

var DefaultChargeController = ChargeController{}

type ChargeController struct {
}

func (ctr ChargeController) Push(c *gin.Context) (gin.H, error) {
	form := api.GetChargeForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	ctx := c.Request.Context()

	//查询 渠道信息
	scene := service.DefaultBdSceneService.FindByCh(form.Ch)
	if scene.Key == "" || scene.Key == "e" {
		return nil, errno.ErrCommon.WithMessage("渠道查询失败")
	}

	//白名单验证
	ip := c.ClientIP()
	if err := service.DefaultBdSceneService.CheckWhiteList(ip, form.Ch); err != nil {
		app.Logger.Info("校验白名单失败", ip)
		return nil, errno.ErrCommon.WithMessage("非白名单ip:" + ip)
	}

	//校验sign
	if scene.Ch != "lvmiao" {
		if !service.DefaultBdSceneService.CheckSign(form.Mobile, form.OutTradeNo, form.TotalPower, form.Sign, scene) {
			app.Logger.Info("校验sign失败", form)
			return nil, errno.ErrCommon.WithMessage("sign:" + form.Sign + " 验证失败")
		}
	}

	//通过手机号查询用户
	userInfo, exist, err := service.DefaultUserService.GetUser(repository.GetUserBy{
		Mobile: form.Mobile,
		Source: entity.UserSourceMio,
	})
	if err != nil {
		return nil, errno.ErrCommon.WithMessage(err.Error())
	}
	//用户验证
	if !exist {
		return nil, errno.ErrCommon.WithMessage("未找到用户")
	}

	//风险登记验证
	if userInfo.Risk >= 4 {
		fmt.Println("用户风险等级过高 ", form)
		return nil, errno.ErrCommon.WithMessage("账户风险等级过高")
	}

	//查重
	transService := service.NewPointTransactionService(mioContext.NewMioContext(mioContext.WithContext(ctx)))
	typeString := service.DefaultBdSceneService.SceneToType(scene.Ch)

	by, err := transService.FindBy(repository.FindPointTransactionBy{
		OpenId: userInfo.OpenId,
		Type:   string(typeString),
		Note:   form.Ch + "#" + form.OutTradeNo,
	})

	if err != nil {
		return nil, err
	}

	if by.ID != 0 {
		app.Logger.Infof("[%s] 重复提交订单", form.Ch)
		return nil, errno.ErrCommon.WithMessage("重复提交订单")
	}

	//入参保存
	defer trackBehaviorInteraction(trackInteractionParam{
		Tp:     string(typeString),
		Data:   form,
		Ip:     c.ClientIP(),
		UserId: userInfo.ID,
	})

	//回调光环
	go ctr.turnPlatform(userInfo, form)

	totalPower, _ := strconv.ParseFloat(form.TotalPower, 64)

	//加碳量
	typeCarbonStr := service.DefaultBdSceneService.SceneToCarbonType(scene.Ch)
	if typeCarbonStr != "" && totalPower != 0 {
		_, errCarbon := service.NewCarbonTransactionService(mioContext.NewMioContext(mioContext.WithContext(ctx))).Create(api_types.CreateCarbonTransactionDto{
			OpenId:  userInfo.OpenId,
			UserId:  userInfo.ID,
			Type:    typeCarbonStr,
			Value:   totalPower,
			Info:    form.OutTradeNo + "#" + form.Mobile + "#" + form.Ch + "#" + fmt.Sprintf("%f", totalPower) + "#" + form.Sign,
			AdminId: 0,
			Ip:      "",
		})
		if errCarbon != nil {
			fmt.Println("charge 加碳失败", form)
		}
	}

	//查询今日积分总量
	timeStr := time.Now().Format("2006-01-02")
	key := timeStr + scene.Ch + form.Mobile
	cmd := app.Redis.Get(ctx, key)
	lastPoint, _ := strconv.Atoi(cmd.Val())
	if lastPoint >= scene.PointLimit {
		fmt.Println("charge 充电量已达到上限 ", form)
		return nil, nil
	}

	thisPoint := int(totalPower * float64(scene.Override))
	totalPoint := lastPoint + thisPoint
	if totalPoint > scene.PointLimit {
		fmt.Println("charge 充电量限制修正 ", form, thisPoint, lastPoint)
		thisPoint = scene.PointLimit - lastPoint
		totalPoint = scene.PointLimit
	}

	app.Redis.Set(ctx, key, totalPoint, 24*36000*time.Second)

	//加积分
	if thisPoint != 0 {
		pointService := service.NewPointService(mioContext.NewMioContext(mioContext.WithContext(ctx)))
		_, err = pointService.IncUserPoint(srv_types.IncUserPointDTO{
			OpenId:       userInfo.OpenId,
			Type:         typeString,
			ChangePoint:  int64(thisPoint),
			BizId:        util.UUID(),
			AdditionInfo: form.OutTradeNo + "#" + form.Mobile + "#" + form.Ch + "#" + strconv.Itoa(thisPoint) + "#" + form.Sign,
		})
		if err != nil {
			app.Logger.Errorf("[%s]加积分失败: %s; query: [%v]\n", form.Ch, err.Error(), form)
		}
	}

	//抽奖
	go func(ctx context.Context, platformKey string, userId int64, power float64) {
		err = ctr.luckyDraw(ctx, scene.Ch, userInfo.ID, totalPower)
		if err != nil {
			app.Logger.Errorf("[%s][%s][%s]错误:%s", "charge", "luckyDraw", scene.Ch, err.Error())
		}
	}(ctx, scene.Ch, userInfo.ID, totalPower)
	//发券
	go func(ctx context.Context, platformKey string, power float64, userInfo *entity.User) {
		err := ctr.sendCoupon(ctx, scene.Ch, totalPower, userInfo)
		if err != nil {
			app.Logger.Errorf("[%s][%s][%s]错误:%s", "charge", "luckyDraw", scene.Ch, err.Error())
		}
	}(ctx, scene.Ch, totalPower, userInfo)
	return gin.H{}, nil
}

func (ctr ChargeController) luckyDraw(ctx context.Context, platformKey string, userId int64, power float64) error {
	if power < 10.00 || platformKey != "lvmiao" {
		return nil
	}
	rule, err := ctr.checkRule(ctx, "starCharge-luckyDraw")
	if err != nil {
		return err
	}

	rdsKey := fmt.Sprintf("%s:%s", config.RedisKey.StarCharge, "starCharge-luckyDraw")
	_, err = app.Redis.SAdd(ctx, rdsKey, userId).Result()
	app.Redis.ExpireAt(ctx, rdsKey, time.UnixMilli(rule.GetEndTime()))
	if err != nil {
		return err
	}
	return nil
}

func (ctr ChargeController) sendCoupon(ctx context.Context, platformKey string, power float64, userInfo *entity.User) error {
	if power < 10.00 || platformKey != "lvmiao" {
		return nil
	}
	rule, err := ctr.checkRule(ctx, platformKey)
	if err != nil {
		app.Logger.Info(fmt.Printf("[charge][sendCoupon] openId:%s ; checkRule失败:%s\n", userInfo.OpenId, err.Error()))
		return err
	}
	//限制判断
	err = ctr.checkLimitOfDay(ctx, platformKey, userInfo.ID)
	if err != nil {
		app.Logger.Info(fmt.Printf("[charge][sendCoupon] openId:%s ; checkLimitOfDay失败:%s\n", userInfo.OpenId, err.Error()))
		return err
	}

	err = ctr.checkLimitOfPeriod(ctx, platformKey, userInfo.ID, 2, rule.GetEndTime())
	if err != nil {
		app.Logger.Info(fmt.Printf("[charge][sendCoupon] openId:%s ; checkLimitOfPeriod失败:%s\n", userInfo.OpenId, err.Error()))
		return err
	}
	//发券
	starChargeService := star_charge.NewStarChargeService(mioContext.NewMioContext(mioContext.WithContext(ctx)))
	token, err := starChargeService.GetAccessToken()
	if err != nil {
		app.Logger.Info(fmt.Printf("[charge][sendCoupon] openId:%s ; 获取token失败:%s\n", userInfo.OpenId, err.Error()))
		return err
	}
	if err = starChargeService.SendCoupon(userInfo.OpenId, userInfo.PhoneNumber, starChargeService.ProvideId, token); err != nil {
		app.Logger.Info(fmt.Printf("[charge][sendCoupon] openId:%s ; 发券失败:%s\n", userInfo.OpenId, err.Error()))
		return err
	}

	return nil
}

func (ctr ChargeController) checkRule(ctx context.Context, platformKey string) (*activity.ActivityRule, error) {
	rule, err := app.RpcService.ActivityRpcSrv.ActiveRule(ctx, &activity.ActiveRuleReq{
		Code: platformKey,
	})
	if err != nil {
		return nil, err
	}
	if !rule.GetExist() {
		return nil, errno.ErrRecordNotFound
	}
	startTime := time.UnixMilli(rule.GetActivityRule().GetStartTime())
	endTime := time.UnixMilli(rule.GetActivityRule().GetEndTime())
	if time.Now().Before(startTime) || time.Now().After(endTime) {
		return nil, errno.ErrMisMatchCondition.WithMessage("活动未开始或已失效")
	}
	return rule.GetActivityRule(), nil
}

func (ctr ChargeController) checkLimitOfDay(ctx context.Context, platformKey string, userId int64) error {
	rdsKey := fmt.Sprintf("%s:%s:%s:%d", config.RedisKey.PeriodLimit, platformKey, time.Now().Format("20060102"), userId)
	periodLimit := limit.NewPeriodLimit(int(time.Hour.Seconds())*24, 1, app.Redis, rdsKey, limit.PeriodAlign())
	res1, err := periodLimit.TakeCtx(ctx, "")
	if err != nil {
		return err
	}
	if res1 != 1 && res1 != 2 {
		return errno.ErrCommon.WithMessage("每日上限1次")
	}
	return nil
}

func (ctr ChargeController) checkLimitOfPeriod(ctx context.Context, platformKey string, userId int64, frequency int, endTime int64) error {
	rdsKey := fmt.Sprintf("%s:%s:%d", config.RedisKey.PeriodLimit, platformKey, userId)
	s := endTime - time.Now().UnixMilli()
	periodLimit2 := limit.NewPeriodLimit(int(s), frequency, app.Redis, rdsKey)
	res2, err := periodLimit2.TakeCtx(ctx, "")
	if err != nil {
		return err
	}
	if res2 != 1 && res2 != 2 {
		return errno.ErrCommon.WithMessage(fmt.Sprintf("活动内上限%d次", frequency))
	}
	return nil
}

func (ctr ChargeController) turnPlatform(user *entity.User, form api.GetChargeForm) {
	sceneUser := repository.DefaultBdSceneUserRepository.FindPlatformUser(user.OpenId, "ccring")
	if sceneUser.ID != 0 && sceneUser.PlatformKey == "ccring" {
		ccringScene := service.DefaultBdSceneService.FindByCh("ccring")
		if ccringScene.ID == 0 {
			app.Logger.Info("ccring 渠道查询失败")
			return
		}
		point, _ := strconv.ParseFloat(form.TotalPower, 64)
		ccRingService := ccring.NewCCRingService("dsaflsdkfjxcmvoxiu123moicuvhoi123", ccringScene.Domain, "/api/cc-ring/external/ev-charge",
			ccring.WithCCRingOrderNum(form.OutTradeNo),
			ccring.WithCCRingMemberId(sceneUser.PlatformUserId),
			ccring.WithCCRingDegreeOfCharge(point),
		)
		//记录
		one := repository.DefaultBdSceneCallbackRepository.FindOne(repository.GetSceneCallback{
			PlatformKey:    sceneUser.PlatformKey,
			PlatformUserId: sceneUser.PlatformUserId,
			OpenId:         sceneUser.OpenId,
			BizId:          form.OutTradeNo,
			SourceKey:      "star_charge",
		})
		if one.ID == 0 {
			body, _ := ccRingService.CallBack()
			err := repository.DefaultBdSceneCallbackRepository.Save(entity.BdSceneCallback{
				PlatformKey:    sceneUser.PlatformKey,
				PlatformUserId: sceneUser.PlatformUserId,
				OpenId:         sceneUser.OpenId,
				BizId:          form.OutTradeNo,
				SourceKey:      "star_charge",
				Body:           body,
				CreatedAt:      time.Now(),
			})
			if err != nil {
				return
			}
		}
	}
}

//充电
func (ctr ChargeController) Ykc(c *gin.Context) (gin.H, error) {
	form := YkcReq{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	ch := "ykc"
	ctx := mioContext.NewMioContext(mioContext.WithContext(c.Request.Context()))

	//查询 渠道信息
	scene := service.DefaultBdSceneService.FindByCh(ch)
	if scene.Key == "" || scene.Key == "e" {
		return nil, errno.ErrCommon.WithMessage("渠道查询失败")
	}

	//查询用户
	userId, err := strconv.ParseInt(form.ExternalUserId, 10, 64)
	if err != nil {
		return nil, err
	}
	userInfo, b, err := service.DefaultUserService.GetUserByID(userId)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, errno.ErrUserNotFound
	}

	//查重
	var bizId, tp string
	bizId = form.TradeSeq
	tp = string(entity.POINT_YKC)
	by, err := app.RpcService.PointRpcSrv.FindPointTransaction(ctx.Context, &point.FindPointTransactionReq{
		BizId: &bizId,
		Type:  &tp,
	})
	if err != nil {
		return nil, err
	}

	if by.Exist {
		app.Logger.Info("云快充 订单重复", form)
		return nil, errno.ErrCommon.WithMessage("重复订单")
	}

	//入参保存
	defer trackBehaviorInteraction(trackInteractionParam{
		Tp:     tp,
		Data:   form,
		Ip:     c.ClientIP(),
		UserId: userId,
	})

	info, _ := json.Marshal(form)
	//加减碳量
	typeCarbonStr := service.DefaultBdSceneService.SceneToCarbonType(scene.Ch)
	if typeCarbonStr != "" && form.ChargedPower != 0 {
		_, errCarbon := service.NewCarbonTransactionService(ctx).Create(api_types.CreateCarbonTransactionDto{
			OpenId:  userInfo.OpenId,
			UserId:  userInfo.ID,
			Type:    entity.CARBON_YKC,
			Value:   form.ChargedPower,
			Info:    string(info),
			AdminId: 0,
			Ip:      "",
		})
		if errCarbon != nil {
			app.Logger.Errorf(fmt.Sprintf("[ykc]加碳失败: %s; query: %v", err.Error(), fmt.Sprint(info)))
		}
	}

	//查询今日积分总量
	timeStr := time.Now().Format("2006-01-02")
	key := timeStr + scene.Ch + form.ExternalUserId
	cmd := app.Redis.Get(ctx, key)

	lastPoint, _ := strconv.Atoi(cmd.Val())
	if lastPoint >= scene.PointLimit {
		return nil, nil
	}

	thisPoint := int(form.ChargedPower * float64(scene.Override))
	totalPoint := lastPoint + thisPoint

	if totalPoint > scene.PointLimit {
		fmt.Printf("%s 充电量限制修正 thisPoint:%d, lastPoint:%d", ch, thisPoint, lastPoint)
		thisPoint = scene.PointLimit - lastPoint
		totalPoint = scene.PointLimit
	}

	app.Redis.Set(ctx, key, totalPoint, 24*36000*time.Second)

	if thisPoint != 0 {
		//加积分
		_, err = app.RpcService.PointRpcSrv.IncPoint(ctx.Context, &point.IncPointReq{
			Openid:       userInfo.OpenId,
			Type:         string(entity.POINT_YKC),
			BizId:        form.TradeSeq,
			BizName:      "云快充订单同步",
			ChangePoint:  uint64(thisPoint),
			AdditionInfo: string(info),
		})
		if err != nil {
			app.Logger.Errorf(fmt.Sprintf("[ykc]加积分失败: %s; query: %v", err.Error(), fmt.Sprint(info)))
		}
	}

	return gin.H{}, nil
}
