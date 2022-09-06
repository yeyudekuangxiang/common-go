package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	"mio/pkg/errno"
	"mio/pkg/oola"
	"strconv"
	"strings"
	"time"
)

var DefaultRecycleController = RecycleController{}

type RecycleController struct {
}

func (ctr RecycleController) OolaOrderSync(c *gin.Context) (gin.H, error) {
	// type oola
	form := RecyclePushForm{}
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
	RecycleService := service.NewRecycleService(ctx)
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
		fmt.Println("charge 未找到用户 ", form)
		return nil, errno.ErrUserNotFound
	}

	//避开重放
	if !util.DefaultLock.Lock(form.Type+form.OrderNo, 24*3600*30*time.Second) {
		fmt.Println("charge 重复提交订单", form)
		app.Logger.Info("charge 重复提交订单", form)
		return nil, errors.New("重复提交订单")
	}
	//if err = RecycleService.CheckOrder(userInfo.OpenId, form.OrderNo); err != nil {
	//	return nil, err
	//}

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
	//本次可得减碳量 todo
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
		AdditionInfo: form.OrderNo + "#" + strconv.FormatInt(currCo2, 10) + "#" + form.ClientId,
		Note:         form.OrderNo,
	})
	if err != nil {
		fmt.Println("oola 旧物回收 加积分失败 ", form)
	}
	return gin.H{}, nil
}

func (ctr RecycleController) GetOolaKey(c *gin.Context) (gin.H, error) {
	//查询 渠道信息
	scene := service.DefaultBdSceneService.FindByCh("oola")
	if scene.Key == "" || scene.Key == "e" {
		return nil, errors.New("渠道查询失败")
	}
	userInfo := apiutil.GetAuthUser(c)
	oolaPkg := oola.NewOola(context.NewMioContext(), scene.AppId, userInfo.OpenId, scene.Domain, app.Redis)
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
	form := RecycleFmyForm{}
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

	dst := service.FmySignParams{}
	err := util.MapTo(&form, &dst)
	if err != nil {
		return nil, err
	}

	//new RecycleService
	ctx := context.NewMioContext()
	RecycleService := service.NewRecycleService(ctx)
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

	//避开重放
	if err = RecycleService.CheckOrder(userInfo.OpenId, form.Data.OrderSn); err != nil {
		return nil, err
	}

	//默认只有衣物
	typeName := entity.POINT_FMY_RECYCLING_CLOTHING
	typeText := RecycleService.GetText(typeName)
	//查询今日该类型获取积分次数
	err = RecycleService.CheckLimit(userInfo.OpenId, typeText)
	if err != nil {
		return nil, err
	}

	//本次可得积分
	currPoint, _ := RecycleService.GetPoint(typeText, form.Data.Weight)
	//本次可得减碳量 todo
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
		AdditionInfo: form.Data.OrderSn + "#" + strconv.FormatInt(currCo2, 10) + "#" + form.Data.Phone,
		Note:         form.Data.OrderSn,
	})
	if err != nil {
		fmt.Println("fmy 旧物回收 加积分失败 ", form)
	}
	return gin.H{}, nil
}
