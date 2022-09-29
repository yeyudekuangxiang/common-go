package open

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	platformService "mio/internal/pkg/service/platform"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	"mio/pkg/errno"
	"mio/pkg/platform"
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
	if form.Type != "1" {
		return nil, errors.New("非回收订单")
	}

	//查询 渠道信息
	scene := service.DefaultBdSceneService.FindByCh("oola")
	if scene.Key == "" || scene.Key == "e" {
		app.Logger.Info("渠道查询失败", form)
		return nil, errors.New("渠道查询失败")
	}

	dst := make(map[string]interface{}, 0)
	err := util.MapTo(&form, &dst)
	if err != nil {
		return nil, err
	}

	//new RecycleService
	ctx := context.NewMioContext()
	RecycleService := platformService.NewRecycleService(ctx)
	TransActionLimitService := service.NewPointTransactionCountLimitService(ctx)
	PointService := service.NewPointService(ctx)
	//校验sign
	if err := RecycleService.CheckSign(dst, scene.Key); err != nil {
		app.Logger.Info("校验sign失败", form)
		return nil, errors.New("sign:" + form.Sign + " 验证失败")
	}
	//通过openid查询用户
	userInfo, _ := service.DefaultUserService.GetUserByOpenId(form.ClientId)
	if userInfo.ID == 0 {
		fmt.Println("oola 未找到用户 ", form)
		return nil, errno.ErrUserNotFound
	}

	//检查重复订单
	if err = RecycleService.CheckOrder(userInfo.OpenId, scene.Ch+"#"+form.OrderNo); err != nil {
		return nil, err
	}

	//回调光环
	go ctr.turnPlatform(userInfo, form)

	//匹配大类型
	typeName := RecycleService.GetType(form.ProductCategoryName)
	if typeName == "" {
		return nil, errors.New("未识别回收分类")
	}

	//查询今日该类型获取积分次数
	err = RecycleService.CheckLimit(userInfo.OpenId, form.Name)
	if err != nil {
		return nil, err
	}

	//本次可得积分
	currPoint, _ := RecycleService.GetPoint(form.Name, form.Qua)
	//本次可得减碳量
	currCo2, _ := RecycleService.GetCo2(form.Name, form.Qua)
	//本月可获得积分上限
	monthPoint, _ := RecycleService.GetMaxPointByMonth(typeName)

	//最终本次回收可获得积分
	point, err := TransActionLimitService.CheckMaxPointByMonth(typeName, userInfo.OpenId, currPoint, monthPoint)
	if err != nil {
		return nil, err
	}
	//加积分
	_, err = PointService.IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       userInfo.OpenId,
		Type:         typeName,
		ChangePoint:  point,
		BizId:        util.UUID(),
		AdditionInfo: form.OrderNo + "#" + strconv.FormatFloat(currCo2, 'f', 2, 64) + "#" + strconv.FormatInt(currPoint, 10) + "#" + form.ClientId,
		Note:         scene.Ch + "#" + form.OrderNo,
	})
	if err != nil {
		fmt.Println("oola 旧物回收 加积分失败 ", form)
	}
	carbonString := fmt.Sprintf("%f", currCo2)

	//发碳量
	carbon, _ := service.NewCarbonTransactionService(context.NewMioContext()).CreateV2(api_types.CreateCarbonTransactionDto{
		OpenId:  userInfo.OpenId,
		UserId:  userInfo.ID,
		Type:    entity.CarbonTransactionType(typeName),
		Value:   currCo2,
		Info:    form.OrderNo + "#" + carbonString + "#" + form.ClientId,
		AdminId: 0,
		Ip:      "",
	})
	println(carbon)

	return gin.H{}, nil
}

func (ctr RecycleController) GetOolaKey(c *gin.Context) (gin.H, error) {
	//查询 渠道信息
	scene := service.DefaultBdSceneService.FindByCh("oola")
	if scene.Key == "" || scene.Key == "e" {
		return nil, errors.New("渠道查询失败")
	}
	userInfo := apiutil.GetAuthUser(c)
	oolaPkg := platform.NewOola(context.NewMioContext(), scene.AppId, userInfo.OpenId, scene.Domain, app.Redis)
	oolaPkg.WithHeadImgUrl(userInfo.AvatarUrl)
	oolaPkg.WithUserName(userInfo.Nickname)
	oolaPkg.WithPhone(userInfo.PhoneNumber)
	channelCode, LoginKey, err := oolaPkg.GetToken()
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
	form := api.RecycleFmyForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	if strings.ToUpper(form.Data.Status) != "COMPLETE" {
		return nil, errors.New("订单未完成")
	}
	//查询 渠道信息
	scene := service.DefaultBdSceneService.FindByCh("fmy")
	if scene.Key == "" || scene.Key == "e" {
		app.Logger.Info("渠道查询失败", form)
		return nil, errors.New("渠道查询失败")
	}

	dst := platformService.FmySignParams{}
	err := util.MapTo(&form, &dst)
	if err != nil {
		return nil, err
	}

	//new RecycleService
	ctx := context.NewMioContext()
	RecycleService := platformService.NewRecycleService(ctx)
	TransActionLimitService := service.NewPointTransactionCountLimitService(ctx)
	PointService := service.NewPointService(ctx)
	//校验sign
	if err := RecycleService.CheckFmySign(dst, scene.AppId, scene.Key); err != nil {
		app.Logger.Info("校验sign失败", form)
		return nil, errors.New("sign:" + form.Sign + " 验证失败")
	}

	//通过phone_number查询用户
	userInfo, _ := service.DefaultUserService.GetUserBy(repository.GetUserBy{
		Source: entity.UserSourceMio,
		Mobile: form.Data.Phone,
	})

	if userInfo.ID == 0 {
		fmt.Println("charge 未找到用户 ", form)
		return nil, errno.ErrUserNotFound
	}

	//默认只有衣物
	typeName := entity.POINT_FMY_RECYCLING_CLOTHING
	typeText := RecycleService.GetText(typeName)

	//幂等 检查重复订单
	if err = RecycleService.CheckOrder(userInfo.OpenId, scene.Ch+"#"+form.Data.OrderSn); err != nil {
		return nil, err
	}

	//查询今日该类型获取积分次数
	err = RecycleService.CheckLimit(userInfo.OpenId, typeText)
	if err != nil {
		return nil, err
	}

	//本次可得积分
	currPoint, _ := RecycleService.GetPoint(typeText, form.Data.Weight)
	//本次可得减碳量
	currCo2, _ := RecycleService.GetCo2(typeText, form.Data.Weight)
	//本月可获得积分上限
	monthPoint, _ := RecycleService.GetMaxPointByMonth(typeName)

	//最终本次回收可获得积分
	point, err := TransActionLimitService.CheckMaxPointByMonth(typeName, userInfo.OpenId, currPoint, monthPoint)
	if err != nil {
		return nil, err
	}
	//加积分
	_, err = PointService.IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       userInfo.OpenId,
		Type:         entity.POINT_FMY_RECYCLING_CLOTHING,
		ChangePoint:  point,
		BizId:        util.UUID(),
		AdditionInfo: form.Data.OrderSn + "#" + strconv.FormatFloat(currCo2, 'E', -1, 64) + "#" + strconv.FormatInt(currPoint, 10) + "#" + form.Data.Phone,
		Note:         scene.Ch + "#" + form.Data.OrderSn,
	})

	carbonString := fmt.Sprintf("%f", currCo2)

	//发碳量
	_, _ = service.NewCarbonTransactionService(context.NewMioContext()).CreateV2(api_types.CreateCarbonTransactionDto{
		OpenId:  userInfo.OpenId,
		UserId:  userInfo.ID,
		Type:    entity.CarbonTransactionType(typeName),
		Value:   currCo2,
		Info:    form.Data.OrderSn + "#" + carbonString + "#" + form.Data.Phone,
		AdminId: 0,
		Ip:      "",
	})

	if err != nil {
		fmt.Println("fmy 旧物回收 加积分失败 ", form)
	}
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
		ccRingService := platformService.NewCCRingService("dsaflsdkfjxcmvoxiu123moicuvhoi123", ccringScene.Domain, "/api/cc-ring/external/recycle",
			platformService.WithCCRingOrderNum(form.OrderNo),
			platformService.WithCCRingMemberId(sceneUser.PlatformUserId),
			platformService.WithCCRingProductCategoryName(form.ProductCategoryName),
			platformService.WithCCRingName(form.Name),
			platformService.WithCCRingQua(form.Qua),
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
