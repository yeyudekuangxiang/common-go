package open

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"mio/internal/app/mp2c/controller/api"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/platform"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	"strconv"
	"time"
)

var DefaultChargeController = ChargeController{}

type ChargeController struct {
}

func (ctr ChargeController) Push(c *gin.Context) (gin.H, error) {
	form := api.GetChargeForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		app.Logger.Errorf("charge/push 参数错误: %s", form)
		return nil, err
	}
	ctx := context.NewMioContext()
	//查询 渠道信息
	scene := service.DefaultBdSceneService.FindByCh(form.Ch)
	if scene.Key == "" || scene.Key == "e" {
		return nil, errors.New("渠道查询失败")
	}
	//白名单验证
	ip := c.ClientIP()
	if err := service.DefaultBdSceneService.CheckWhiteList(ip, form.Ch); err != nil {
		app.Logger.Info("校验白名单失败", ip)
		return nil, errors.New("非白名单ip:" + ip)
	}

	//校验sign
	if scene.Ch != "lvmiao" {
		if !service.DefaultBdSceneService.CheckSign(form.Mobile, form.OutTradeNo, form.TotalPower, form.Sign, scene) {
			app.Logger.Info("校验sign失败", form)
			return nil, errors.New("sign:" + form.Sign + " 验证失败")
		}
	}

	//避开重放
	if !util.DefaultLock.Lock(form.Ch+form.OutTradeNo, 24*3600*30*time.Second) {
		fmt.Println("charge 重复提交订单", form)
		app.Logger.Info("charge 重复提交订单", form)
		return nil, errors.New("重复提交订单")
	}

	//通过手机号查询用户
	userInfo, _ := service.DefaultUserService.GetUserBy(repository.GetUserBy{
		Mobile: form.Mobile,
		Source: entity.UserSourceMio,
	})

	if userInfo.ID <= 0 {
		fmt.Println("charge 未找到用户 ", form)
		return nil, errors.New("未找到用户")
	}

	//风险登记验证
	if userInfo.Risk >= 2 {
		fmt.Println("用户风险等级过高 ", form)
		return nil, errors.New("账户风险等级过高")
	}

	//查询今日积分总量
	timeStr := time.Now().Format("2006-01-02")
	key := timeStr + scene.Ch + form.Mobile
	cmd := app.Redis.Get(ctx, key)

	lastPoint, _ := strconv.Atoi(cmd.Val())
	thisPoint0, _ := strconv.ParseFloat(form.TotalPower, 64)

	thisPoint := int(thisPoint0 * float64(scene.Override))
	totalPoint := lastPoint + thisPoint
	if lastPoint >= scene.PointLimit {
		fmt.Println("charge 充电量已达到上限 ", form)
	} else {
		if totalPoint > scene.PointLimit {
			fmt.Println("charge 充电量限制修正 ", form, thisPoint, lastPoint)
			thisPoint = scene.PointLimit - lastPoint
			totalPoint = scene.PointLimit
		}
	}

	app.Redis.Set(ctx, key, totalPoint, 24*36000*time.Second)

	//加积分
	typeString := service.DefaultBdSceneService.SceneToType(scene.Ch)
	pointService := service.NewPointService(ctx)
	_, err := pointService.IncUserPoint(srv_types.IncUserPointDTO{
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
	////绿喵回调第三方
	//ccRingService := platform.NewCCRingService()
	//go ccRingService.CallBack(userInfo, thisPoint0, scene.Ch, scene)
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

//调用星星充电
func (ctr ChargeController) sendCoupon(ctx *context.MioContext, platformKey string, point int64, userInfo *entity.User) {
	if app.Redis.Exists(ctx, platformKey+"_"+"ChargeException").Val() == 0 && point > 0 {
		fmt.Println("星星充电 发券start")
		startTime, _ := time.Parse("2006-01-02", "2022-08-22")
		endTime, _ := time.Parse("2006-01-02", "2022-08-31")
		if platformKey == "lvmiao" && time.Now().After(startTime) && time.Now().Before(endTime) {
			starChargeService := platform.NewStarChargeService(ctx)
			token, err := starChargeService.GetAccessToken()
			if err != nil {
				fmt.Printf("星星充电 获取token失败:%s\n", err.Error())
				app.Logger.Info(fmt.Printf("星星充电 openId:%s ; 获取token失败:%s\n", userInfo.OpenId, err.Error()))
			}
			//限制一次
			if err = starChargeService.CheckChargeLimit(userInfo.OpenId); err != nil {
				fmt.Printf("星星充电 检查次数限制:%s\n", err.Error())
				app.Logger.Info(fmt.Printf("星星充电 openId:%s ; 检查次数限制:%s\n", userInfo.OpenId, err.Error()))
			}
			//send coupon
			if err = starChargeService.SendCoupon(userInfo.OpenId, userInfo.PhoneNumber, starChargeService.ProvideId, token); err != nil {
				fmt.Printf("星星充电 发券失败:%s\n", err.Error())
				app.Logger.Info(fmt.Printf("星星充电 openId:%s ; 发券失败:%s\n", userInfo.OpenId, err.Error()))
			}
		}
	}
}
