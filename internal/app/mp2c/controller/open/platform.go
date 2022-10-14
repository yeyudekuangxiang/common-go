package open

import (
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
	"mio/internal/pkg/service/platform/jhx"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	"mio/pkg/errno"
	platformUtil "mio/pkg/platform"
	"strconv"
	"strings"
	"time"
)

var DefaultPlatformController = PlatformController{}

type PlatformController struct {
}

func (receiver PlatformController) BindPlatformUser(ctx *gin.Context) (gin.H, error) {
	//接收参数 platformKey, phone
	form := bindPlatform{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)
	//查询渠道号
	scene := service.DefaultBdSceneService.FindByCh(form.PlatformKey)
	if scene.Key == "" || scene.Key == "e" {
		app.Logger.Info("渠道查询失败", form)
		return nil, errno.ErrCommon.WithMessage("渠道查询失败")
	}

	if user.ID == 0 {
		return nil, errno.ErrCommon.WithMessage("用户未登陆")
	}

	//绑定
	sceneUser, err := service.DefaultBdSceneUserService.Bind(user, *scene, form.MemberId)
	if err != nil && err != errno.ErrChannelExisting {
		return nil, err
	}
	//绑定回调
	if scene.Ch == "jinhuaxing" && err != errno.ErrChannelExisting {
		err = jhx.NewJhxService(context.NewMioContext()).BindSuccess(sceneUser.Phone, "1")
		if err != nil {
			app.Logger.Errorf("callback %s error:%s", scene.Ch, err.Error())
			return nil, err
		}
	}
	//返回
	return gin.H{
		"memberId":     sceneUser.PlatformUserId,
		"lvmiaoUserId": sceneUser.OpenId,
	}, nil
}

func (receiver PlatformController) SyncPoint(ctx *gin.Context) (gin.H, error) {
	form := platformForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	dst := make(map[string]interface{}, 0)
	err := util.MapTo(&form, &dst)
	if err != nil {
		return nil, err
	}

	//查询渠道号
	scene := service.DefaultBdSceneService.FindByCh(form.PlatformKey)
	if scene.Key == "" || scene.Key == "e" {
		app.Logger.Info("渠道查询失败", form)
		return nil, errno.ErrCommon.WithMessage("渠道查询失败")
	}
	delete(dst, "sign")
	//check sign
	if err := platformUtil.CheckSign(form.Sign, dst, form.PlatformKey, ";"); err != nil {
		app.Logger.Errorf("校验sign失败: %s", err.Error())
		return nil, err
	}

	//check user
	user, _ := service.DefaultUserService.GetUserBy(repository.GetUserBy{Mobile: form.Mobile, Source: entity.UserSourceMio})
	if user.ID == 0 {
		return nil, errno.ErrCommon.WithMessage("用户不存在")
	}

	method := scene.Ch
	if form.Method != "" {
		method = strings.ToLower(method) + "_" + strings.ToLower(form.Method)
	}
	if _, ok := entity.PlatformMethodMap[method]; !ok {
		return nil, errno.ErrCommon.WithMessage("未找到匹配方法")
	}
	t := entity.PlatformMethodMap[method]

	value := entity.PointCollectValueMap[t]

	_, err = service.NewPointService(context.NewMioContext()).IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:      user.OpenId,
		Type:        t,
		BizId:       util.UUID(),
		ChangePoint: int64(value),
		AdminId:     0,
		Note:        t.Text(),
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// PrePoint 预加积分
func (receiver PlatformController) PrePoint(c *gin.Context) (gin.H, error) {
	form := api.PrePointRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		app.Logger.Errorf("参数错误: %s", form)
		return nil, err
	}
	ctx := context.NewMioContext()
	//查询 渠道信息
	scene := service.DefaultBdSceneService.FindByCh(form.PlatformKey)
	if scene.Key == "" || scene.Key == "e" {
		return nil, errno.ErrCommon.WithMessage("渠道查询失败")
	}
	//白名单验证
	ip := c.ClientIP()
	if err := service.DefaultBdSceneService.CheckWhiteList(ip, form.PlatformKey); err != nil {
		app.Logger.Info("校验白名单失败", ip)
		return nil, errno.ErrCommon.WithMessage("非白名单ip:" + ip)
	}

	//校验sign
	params := make(map[string]string, 0)
	err := util.MapTo(&form, &params)
	if err != nil {
		return nil, err
	}
	sign := form.Sign
	delete(params, "sign")
	if !service.DefaultBdSceneService.CheckPreSign(scene.Key, sign, params) {
		app.Logger.Info("校验sign失败", form)
		return nil, errno.ErrCommon.WithMessage("sign:" + form.Sign + " 验证失败")
	}

	sceneUser := repository.DefaultBdSceneUserRepository.FindOne(repository.GetSceneUserOne{
		PlatformKey:    form.PlatformKey,
		PlatformUserId: form.MemberId,
	})

	//查重
	transService := service.NewPointTransactionService(ctx)
	typeString := service.DefaultBdSceneService.SceneToType(scene.Ch)

	by, err := transService.FindBy(repository.FindPointTransactionBy{
		Type: string(typeString),
		Note: form.PlatformKey + "#" + form.TradeNo,
	})

	if err != nil {
		return nil, err
	}

	if by.ID != 0 {
		fmt.Println("charge 重复提交订单", form)
		app.Logger.Info("charge 重复提交订单", form)
		return nil, errno.ErrCommon.WithMessage("重复提交订单")
	}

	//预加积分
	fromString, _ := decimal.NewFromString(params["amount"])
	point := fromString.Mul(decimal.NewFromInt(int64(scene.Override))).Round(0).String()
	err = repository.DefaultBdScenePrePointRepository.Create(&entity.BdScenePrePoint{
		PlatformKey:    form.PlatformKey,
		PlatformUserId: form.MemberId,
		Point:          point,
		OpenId:         sceneUser.OpenId,
		Status:         1,
		Mobile:         sceneUser.Phone,
		Tradeno:        form.TradeNo,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

// GetPrePointList 获取预加积分列表
func (receiver PlatformController) GetPrePointList(c *gin.Context) (gin.H, error) {
	return nil, nil
}

func (receiver PlatformController) CollectPrePoint(c *gin.Context) (gin.H, error) {
	form := api.CollectPrePoint{}
	if err := apiutil.BindForm(c, &form); err != nil {
		app.Logger.Errorf("参数错误: %s", form)
		return nil, err
	}

	//查询 渠道信息
	scene := service.DefaultBdSceneService.FindByCh(form.PlatformKey)
	if scene.Key == "" || scene.Key == "e" {
		return nil, errno.ErrCommon.WithMessage("渠道查询失败")
	}

	//白名单验证
	ip := c.ClientIP()
	if err := service.DefaultBdSceneService.CheckWhiteList(ip, form.PlatformKey); err != nil {
		app.Logger.Info("校验白名单失败", ip)
		return nil, errno.ErrCommon.WithMessage("非白名单ip:" + ip)
	}

	//校验sign
	params := make(map[string]interface{}, 0)
	err := util.MapTo(&form, &params)
	if err != nil {
		return nil, err
	}

	sign := form.Sign
	delete(params, "sign")

	if err = platformUtil.CheckSign(sign, params, scene.Key, "&"); err != nil {
		app.Logger.Info("校验sign失败", form)
		return nil, errno.ErrCommon.WithMessage(err.Error())
	}

	sceneUser := repository.DefaultBdSceneUserRepository.FindOne(repository.GetSceneUserOne{
		PlatformKey:    form.PlatformKey,
		PlatformUserId: form.MemberId,
	})

	if sceneUser.OpenId == "" {
		app.Logger.Infof("未找到绑定关系 platformKey:%s; memberId: %s", form.PlatformKey, form.MemberId)
		return nil, errno.ErrBindRecordNotFound
	}

	userInfo, err := service.DefaultUserService.GetUserByOpenId(sceneUser.OpenId)
	if err != nil {
		return nil, errno.ErrUserNotFound
	}

	//获取pre_point数据 one limit
	id, _ := strconv.ParseInt(form.PrePointId, 10, 64)
	one, err := repository.DefaultBdScenePrePointRepository.FindOne(repository.GetScenePrePoint{
		PlatformKey:    sceneUser.PlatformKey,
		PlatformUserId: sceneUser.PlatformUserId,
		OpenId:         sceneUser.OpenId,
		Id:             id,
		Status:         1,
	})

	if err != nil {
		return nil, errno.ErrRecordNotFound
	}

	//检查上限
	ctx := context.NewMioContext()

	var isHalf bool
	var halfPoint int64

	timeStr := time.Now().Format("2006-01-02")
	key := timeStr + ":prePoint:" + scene.Ch + sceneUser.PlatformUserId + sceneUser.Phone

	lastPoint, _ := strconv.ParseInt(app.Redis.Get(ctx, key).Val(), 10, 64)
	incPoint, _ := strconv.ParseInt(one.Point, 10, 64)

	totalPoint := lastPoint + incPoint

	if lastPoint >= int64(scene.PrePointLimit) {
		return nil, errno.ErrCommon.WithMessage("今日获取积分已达到上限")
	}

	if totalPoint > int64(scene.PrePointLimit) {
		p := incPoint
		incPoint = int64(scene.PrePointLimit) - lastPoint
		totalPoint = int64(scene.PrePointLimit)
		isHalf = true
		halfPoint = p - incPoint
	}

	app.Redis.Set(ctx, key, totalPoint, 24*time.Hour)

	//积分
	point, err := service.NewPointService(context.NewMioContext()).IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:      sceneUser.OpenId,
		Type:        entity.POINT_JHX,
		BizId:       util.UUID(),
		ChangePoint: incPoint,
		AdminId:     0,
		Note:        form.PlatformKey + "#" + one.Tradeno,
	})

	if err != nil {
		return nil, err
	}

	//更新pre_point对应数据
	one.Status = 2
	one.UpdatedAt = time.Now()
	if isHalf {
		one.Status = 1
		one.Point = strconv.FormatInt(halfPoint, 10)
	}

	err = repository.DefaultBdScenePrePointRepository.Save(&one)
	if err != nil {
		return nil, err
	}

	//减碳量
	fromString, _ := decimal.NewFromString(one.Point)
	amount, _ := fromString.Div(decimal.NewFromInt(int64(scene.Override))).Float64()
	typeCarbonStr := service.DefaultBdSceneService.SceneToCarbonType(scene.Ch)
	if typeCarbonStr != "" {
		_, err = service.NewCarbonTransactionService(context.NewMioContext()).Create(api_types.CreateCarbonTransactionDto{
			OpenId: sceneUser.OpenId,
			UserId: userInfo.ID,
			Type:   typeCarbonStr,
			Value:  amount,
			Ip:     ip,
		})
		if err != nil {
			app.Logger.Errorf("预加积分 err:%s", err.Error())
		}

	}

	return gin.H{
		"point":     point,
		"thisPoint": strconv.FormatInt(halfPoint, 10),
	}, nil
}
