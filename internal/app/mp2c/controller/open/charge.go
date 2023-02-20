package open

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/activity/cmd/rpc/activity/activity"
	point2 "gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/point/cmd/rpc/point"
	"mio/internal/app/mp2c/controller/api"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/platform/ccring"
	"mio/internal/pkg/service/platform/star_charge"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
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

	ctx := context.NewMioContext()

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
	transService := service.NewPointTransactionService(ctx)
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
		Tp:   string(typeString),
		Data: form,
		Ip:   c.ClientIP(),
	})

	//查询今日积分总量
	timeStr := time.Now().Format("2006-01-02")
	key := timeStr + scene.Ch + form.Mobile
	cmd := app.Redis.Get(ctx, key)

	lastPoint, _ := strconv.Atoi(cmd.Val())
	thisPoint0, _ := strconv.ParseFloat(form.TotalPower, 64)

	//回调光环
	go ctr.turnPlatform(userInfo, form)

	thisPoint := int(thisPoint0 * float64(scene.Override))
	totalPoint := lastPoint + thisPoint
	if lastPoint >= scene.PointLimit {
		fmt.Println("charge 充电量已达到上限 ", form)
		return nil, nil
		//return nil, errors.New("充电获取积分已达到上限")
	}

	if totalPoint > scene.PointLimit {
		fmt.Println("charge 充电量限制修正 ", form, thisPoint, lastPoint)
		thisPoint = scene.PointLimit - lastPoint
		totalPoint = scene.PointLimit
	}

	app.Redis.Set(ctx, key, totalPoint, 24*36000*time.Second)

	//加积分
	pointService := service.NewPointService(ctx)
	_, err = pointService.IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       userInfo.OpenId,
		Type:         typeString,
		ChangePoint:  int64(thisPoint),
		BizId:        util.UUID(),
		AdditionInfo: form.OutTradeNo + "#" + form.Mobile + "#" + form.Ch + "#" + strconv.Itoa(thisPoint) + "#" + form.Sign,
	})

	if err != nil {
		fmt.Println("charge 加积分失败 ", form)
	}

	//加碳量
	typeCarbonStr := service.DefaultBdSceneService.SceneToCarbonType(scene.Ch)
	if typeCarbonStr != "" {
		pointDec := decimal.NewFromInt(int64(thisPoint))
		electric := pointDec.Div(decimal.NewFromInt(10))
		f, _ := electric.Float64()
		_, errCarbon := service.NewCarbonTransactionService(context.NewMioContext()).Create(api_types.CreateCarbonTransactionDto{
			OpenId:  userInfo.OpenId,
			UserId:  userInfo.ID,
			Type:    typeCarbonStr,
			Value:   f,
			Info:    form.OutTradeNo + "#" + form.Mobile + "#" + form.Ch + "#" + strconv.Itoa(thisPoint) + "#" + form.Sign,
			AdminId: 0,
			Ip:      "",
		})
		if errCarbon != nil {
			fmt.Println("charge 加碳失败", form)
		}
	}

	//发券
	go ctr.sendCoupon(ctx, scene.Ch, int64(thisPoint), userInfo)
	return gin.H{}, nil
}

func (ctr ChargeController) SetException(c *gin.Context) (gin.H, error) {
	form := api.ChangeChargeExceptionForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	ctx := context.NewMioContext()
	err := app.Redis.Set(ctx, form.Ch+"_"+"ChargeException", 1, 0).Err()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ctr ChargeController) DelException(c *gin.Context) (gin.H, error) {
	form := api.ChangeChargeExceptionForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	ctx := context.NewMioContext()
	app.Redis.Del(ctx, form.Ch+"_"+"ChargeException")
	return nil, nil
}

//调用星星充电发券
func (ctr ChargeController) sendCoupon(ctx *context.MioContext, platformKey string, point int64, userInfo *entity.User) {
	if app.Redis.Exists(ctx, platformKey+"_"+"ChargeException").Val() == 0 && point > 0 && platformKey == "lvmiao" {
		rule, err := app.RpcService.ActivityRpcSrv.ActiveRule(ctx.Context, &activity.ActiveRuleReq{
			Code: platformKey,
		})
		if err != nil {
			app.Logger.Info(fmt.Printf("星星充电 openId:[%s]发券失败:%s\n", userInfo.OpenId, err.Error()))
			return
		}
		if !rule.GetExist() {
			app.Logger.Info(fmt.Printf("星星充电 openId:[%s]发券失败:%s\n", userInfo.OpenId, "无有效规则"))
			return
		}
		startTime := time.UnixMilli(rule.GetActivityRule().GetStartTime())
		endTime := time.UnixMilli(rule.GetActivityRule().GetEndTime())
		if time.Now().After(startTime) && time.Now().Before(endTime) {
			starChargeService := star_charge.NewStarChargeService(ctx)
			token, err := starChargeService.GetAccessToken()
			if err != nil {
				app.Logger.Info(fmt.Printf("星星充电 openId:%s ; 获取token失败:%s\n", userInfo.OpenId, err.Error()))
				return
			}
			//限制一次
			if err = starChargeService.CheckChargeLimit(userInfo.OpenId, endTime); err != nil {
				app.Logger.Info(fmt.Printf("星星充电 openId:%s ; 检查次数限制:%s\n", userInfo.OpenId, err.Error()))
				return
			}
			//send coupon
			if err = starChargeService.SendCoupon(userInfo.OpenId, userInfo.PhoneNumber, starChargeService.ProvideId, token); err != nil {
				app.Logger.Info(fmt.Printf("星星充电 openId:%s ; 发券失败:%s\n", userInfo.OpenId, err.Error()))
				return
			}
			return
		}
		return
	}
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
	ctx := context.NewMioContext()

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
	tp = string(entity.Point_YKC)
	by, err := app.RpcService.PointRpcSrv.FindPointTransaction(ctx.Context, &point2.FindPointTransactionReq{
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
		Tp:   tp,
		Data: form,
		Ip:   c.ClientIP(),
	})

	//查询今日积分总量
	timeStr := time.Now().Format("2006-01-02")
	key := timeStr + scene.Ch + form.ExternalUserId
	cmd := app.Redis.Get(ctx, key)

	lastPoint, _ := strconv.Atoi(cmd.Val())
	thisPoint0 := form.ChargedPower

	thisPoint := int(thisPoint0 * float64(scene.Override))
	totalPoint := lastPoint + thisPoint
	if lastPoint >= scene.PointLimit {
		return nil, nil
	}

	if totalPoint > scene.PointLimit {
		fmt.Printf("%s 充电量限制修正 thisPoint:%d, lastPoint:%d", ch, thisPoint, lastPoint)
		thisPoint = scene.PointLimit - lastPoint
		totalPoint = scene.PointLimit
	}

	app.Redis.Set(ctx, key, totalPoint, 24*36000*time.Second)
	marshal, err := json.Marshal(form)
	if err != nil {
		app.Logger.Errorf("云快充 info:%s", err.Error())
	}
	//加积分
	_, err = app.RpcService.PointRpcSrv.IncPoint(ctx.Context, &point2.IncPointReq{
		Openid:       userInfo.OpenId,
		Type:         string(entity.Point_YKC),
		BizId:        form.TradeSeq,
		BizName:      "云快充订单同步",
		ChangePoint:  uint64(thisPoint),
		AdditionInfo: string(marshal),
	})
	if err != nil {
		app.Logger.Errorf("云快充 加积分失败:%s", err.Error())
	}

	//加碳量
	pointDec := decimal.NewFromInt(int64(thisPoint))
	electric := pointDec.Div(decimal.NewFromInt(10))
	f, _ := electric.Float64()
	_, errCarbon := service.NewCarbonTransactionService(context.NewMioContext()).Create(api_types.CreateCarbonTransactionDto{
		OpenId:  userInfo.OpenId,
		UserId:  userInfo.ID,
		Type:    entity.CARBON_ECAR,
		Value:   f,
		Info:    string(marshal),
		AdminId: 0,
		Ip:      "",
	})
	if errCarbon != nil {
		app.Logger.Errorf("云快充 加碳失败:%s", err.Error())
	}

	return gin.H{}, nil
}
