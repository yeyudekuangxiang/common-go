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
	"mio/pkg/errno"
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
	//查询 渠道信息
	scene := service.DefaultBdSceneService.FindByCh(form.Ch)
	if scene.Key == "" || scene.Key == "e" {
		app.Logger.Info("渠道查询失败", form)
		return nil, errors.New("渠道查询失败")
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
	} else {
		//查询今日积分总量
		timeStr := time.Now().Format("2006-01-02")
		key := timeStr + "XingXing" + scene.Ch + form.Mobile
		cmd := app.Redis.Get(c, key)
		lastPoint, _ := strconv.Atoi(cmd.Val())
		thisPoint0, _ := strconv.ParseFloat(form.TotalPower, 64)

		thisPoint := int(thisPoint0 * float64(scene.Override))
		totalPoint := lastPoint + thisPoint
		fmt.Println("charge 累计积分 ", form, totalPoint)
		if lastPoint >= scene.PointLimit {
			fmt.Println("charge 充电量已达到上限 ", form)
		} else {
			if totalPoint > scene.PointLimit {
				fmt.Println("charge 充电量限制修正 ", form, thisPoint, lastPoint)
				thisPoint = scene.PointLimit - lastPoint
				totalPoint = scene.PointLimit
			}
		}
		app.Redis.Set(c, key, totalPoint, 24*36000*time.Second)

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
	}

	return gin.H{}, nil
}

// SendCoupon 发放优惠券
func (ctr ChargeController) SendCoupon(c *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(c)
	if user.PhoneNumber == "" {
		return nil, errno.ErrBind.WithMessage("未绑定手机号")
	}
	XingService := service.XingXingService{
		OperatorID:     "313744932",
		OperatorSecret: "acb93539fc9bg78k",
		SigSecret:      "9af2e7b2d7562ad5",
		DataSecret:     "a2164ada0026ccf7",
		DataSecretIV:   "82c91325e74bef0f",
		Url:            "http://test-evcs.starcharge.com/evcs/starcharge",
	}
	token, err := XingService.GetXingAccessToken(context.NewMioContext())
	if err != nil {
		return nil, err
	}
	fmt.Printf("xing xing accessToken:%s", token)
	err = XingService.SendCount(context.NewMioContext(), user.PhoneNumber, "", token)
	if err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

//快电
func (ctr ChargeController) FastElectric(c *gin.Context) (gin.H, error) {
	form := GetChargeForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	fmt.Println("charge", form)
	//查询 渠道信息
	scene := service.DefaultBdSceneService.FindByCh(form.Ch)
	if scene.Key == "" || scene.Key == "e" {
		app.Logger.Info("渠道查询失败", form)
		return nil, errors.New("渠道查询失败")
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
	} else {
		//查询今日积分总量
		timeStr := time.Now().Format("2006-01-02")
		key := timeStr + "FastElectric" + scene.Ch + form.Mobile
		cmd := app.Redis.Get(c, key)
		lastPoint, _ := strconv.Atoi(cmd.Val())
		thisPoint0, _ := strconv.ParseFloat(form.TotalPower, 64)

		thisPoint := int(thisPoint0 * float64(scene.Override))
		totalPoint := lastPoint + thisPoint
		fmt.Println("charge 累计积分 ", form, totalPoint)
		if lastPoint >= scene.PointLimit {
			fmt.Println("charge 充电量已达到上限 ", form)
		} else {
			if totalPoint > scene.PointLimit {
				fmt.Println("charge 充电量限制修正 ", form, thisPoint, lastPoint)
				thisPoint = scene.PointLimit - lastPoint
				totalPoint = scene.PointLimit
			}
		}
		app.Redis.Set(c, key, totalPoint, 24*36000*time.Second)

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
	}

	return gin.H{}, nil
}
