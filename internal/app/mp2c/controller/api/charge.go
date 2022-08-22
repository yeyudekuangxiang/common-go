package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
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
	form := GetChargeForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	fmt.Println("charge", form)
	ctx := context.NewMioContext()
	//查询 渠道信息
	scene := service.DefaultBdSceneService.FindByCh(form.Ch)
	if scene.Key == "" || scene.Key == "e" {
		app.Logger.Info("渠道查询失败", form)
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
	userInfo, _ := service.DefaultUserService.GetUserBy(repository.GetUserBy{Mobile: form.Mobile})
	if userInfo.ID <= 0 {
		fmt.Println("charge 未找到用户 ", form)
		return nil, errors.New("未找到用户")
	}
	//风险登记验证
	if userInfo.Risk >= 2 {
		fmt.Println("用户风险登记过高 ", form)
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
	fmt.Println("charge 累计积分 ", form, totalPoint)
	if lastPoint >= scene.PointLimit {
		fmt.Println("charge 充电量已达到上线 ", form)
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
	pointService := service.NewPointService(context.NewMioContext())
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
	// todo 发券
	if app.Redis.Exists(ctx, scene.Ch+"Exception").Val() == 0 {
		startTime, _ := time.Parse("20060102", "2022-08-22")
		endTime, _ := time.Parse("20060102", "2022-08-30")
		if scene.Ch == "lvmiao" && time.Now().After(startTime) && time.Now().Before(endTime) {
			starChargeService := service.NewStarChargeService(context.NewMioContext())
			token, err := starChargeService.GetAccessToken()
			if err != nil {
				return nil, err
			}
			//限制一次
			if err = starChargeService.CheckLimit(userInfo.OpenId); err != nil {
				return nil, err
			}
			//send coupon
			if err = starChargeService.SendCoupon(userInfo.OpenId, userInfo.PhoneNumber, starChargeService.ProvideId, token); err != nil {
				return nil, err
			}
		}
	}
	return gin.H{}, nil
}

func (ctr ChargeController) SetException(c *gin.Context) (gin.H, error) {
	form := GetChargeForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	ctx := context.NewMioContext()
	app.Redis.Set(ctx, form.Ch+"Exception", 1, 0)
	return nil, nil
}

func (ctr ChargeController) DelException(c *gin.Context) (gin.H, error) {
	form := GetChargeForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	ctx := context.NewMioContext()
	app.Redis.Del(ctx, form.Ch+"Exception")
	return nil, nil
}
