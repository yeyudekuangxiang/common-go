package open

import (
	"github.com/avast/retry-go"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"mio/config"
	"mio/internal/app/mp2c/controller/api"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/platform/jhx"
	"mio/internal/pkg/service/platform/ytx"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/service/track"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	platformUtil "mio/internal/pkg/util/platform"
	"mio/internal/pkg/util/validator"
	"mio/pkg/errno"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var DefaultPlatformController = PlatformController{}

type PlatformController struct {
}

func (receiver PlatformController) BindPlatformUser(c *gin.Context) (gin.H, error) {
	//接收参数 platformKey, phone
	form := bindPlatform{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)
	//查询渠道号
	scene := service.DefaultBdSceneService.FindByCh(form.PlatformKey)

	if scene.Key == "" || scene.Key == "e" {
		app.Logger.Info("渠道查询失败", form)
		return nil, errno.ErrCommon.WithMessage("第三方绑定 渠道查询失败")
	}

	if user.ID == 0 {
		return nil, errno.ErrCommon.WithMessage("第三方绑定 用户未登陆")
	}

	app.Logger.Infof("第三方绑定 解码前: platformId:%s; openId:%s", form.MemberId, user.OpenId)
	memberId, _ := url.QueryUnescape(form.MemberId)
	app.Logger.Infof("第三方绑定 解码后: platformId:%s; openId:%s", memberId, user.OpenId)

	sceneUser, err := service.DefaultBdSceneUserService.Bind(user, *scene, memberId)
	if err != nil {
		if err != errno.ErrExisting {
			app.Logger.Errorf("第三方绑定 绑定失败:%s; platformId:%s; openId:%s", err.Error(), memberId, user.OpenId)
			return nil, errno.ErrCommon.WithMessage(scene.Ch + "用户绑定失败")
		}
		return gin.H{
			"memberId":     sceneUser.PlatformUserId,
			"lvmiaoUserId": sceneUser.OpenId,
		}, nil
	}

	//绑定回调
	err = retry.Do(
		func() error {
			//绑定回调
			var err error
			//var ch_key string
			if scene.Ch == "jinhuaxing" {
				params := make(map[string]interface{}, 0)
				params["mobile"] = user.PhoneNumber
				params["status"] = "1"
				jhxSvr := jhx.NewJhxService(context.NewMioContext())
				err = jhxSvr.BindSuccess(params)
				//ch_key = "金华行"
			}

			if scene.Ch == "yitongxing" {
				params := make(map[string]interface{}, 0)
				params["memberId"] = sceneUser.PlatformUserId
				params["openId"] = sceneUser.OpenId
				ytxSrv := ytx.NewYtxService(context.NewMioContext(), ytx.WithSecret(scene.Secret2), ytx.WithDomain(scene.Domain2))
				err = ytxSrv.BindSuccess(params)
				//ch_key = "亿通行"
			}
			return err
		},
	)

	if err != nil {
		app.Logger.Errorf("第三方绑定 注册回调失败:%s; platformId:%s; openId:%s", err.Error(), memberId, user.OpenId)
	}

	//返回
	return gin.H{
		"memberId":     sceneUser.PlatformUserId,
		"lvmiaoUserId": sceneUser.OpenId,
	}, nil
}

func (receiver PlatformController) SyncPoint(c *gin.Context) (gin.H, error) {
	form := platformForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
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
	form := setPrePointRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
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
	params := make(map[string]interface{}, 0)
	err := util.MapTo(&form, &params)
	if err != nil {
		return nil, err
	}

	sign := form.Sign
	delete(params, "sign")

	if err = platformUtil.CheckSign(sign, params, scene.Key, "&"); err != nil {
		app.Logger.Info("校验sign失败", form)
		return nil, errno.ErrCommon.WithMessage("sign:" + form.Sign + " 验证失败")
	}

	//查重
	PointTransService := service.NewPointTransactionService(ctx)
	typeString := service.DefaultBdSceneService.SceneToType(scene.Ch)

	by, err := PointTransService.FindBy(repository.FindPointTransactionBy{
		Type: string(typeString),
		Note: form.PlatformKey + "#" + form.TradeNo,
	})

	if err != nil {
		return nil, errno.ErrCommon.WithErr(err)
	}

	if by.ID != 0 {
		app.Logger.Errorf("重复提交订单: %v", form)
		return nil, errno.ErrCommon.WithMessage("重复提交订单")
	}

	var openId, mobile string
	sceneUser := repository.DefaultBdSceneUserRepository.FindOne(repository.GetSceneUserOne{
		PlatformKey:    form.PlatformKey,
		PlatformUserId: form.MemberId,
	})

	if sceneUser.ID != 0 {
		openId = sceneUser.OpenId
		mobile = sceneUser.Phone
	}

	//预加积分
	err = repository.DefaultBdScenePrePointRepository.Create(&entity.BdScenePrePoint{
		PlatformKey:    form.PlatformKey,
		PlatformUserId: form.MemberId,
		Point:          form.Point,
		Status:         1,
		Tradeno:        form.TradeNo,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		OpenId:         openId,
		Mobile:         mobile,
	})

	if err != nil {
		return nil, errno.ErrCommon.WithErr(err)
	}

	return gin.H{}, nil
}

// GetPrePointList 获取预加积分列表
func (receiver PlatformController) GetPrePointList(c *gin.Context) (gin.H, error) {
	form := prePointListRequest{}
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
		return nil, errno.ErrCommon.WithMessage("sign:" + form.Sign + " 验证失败")
	}

	res, total, err := repository.DefaultBdScenePrePointRepository.FindBy(repository.GetScenePrePoint{
		PlatformKey:    form.PlatformKey,
		PlatformUserId: form.MemberId,
		Status:         1,
		StartTime:      time.Now().AddDate(0, 0, -7),
	})

	if err != nil {
		return nil, errno.ErrCommon.WithErr(err)
	}
	var point int64
	sceneUserCondition := repository.GetSceneUserOne{
		PlatformKey:    form.PlatformKey,
		PlatformUserId: form.MemberId,
	}
	sceneUser := service.DefaultBdSceneUserService.FindOne(sceneUserCondition)
	if sceneUser.ID != 0 {
		pointInfo, err := service.NewPointService(context.NewMioContext()).FindByOpenId(sceneUser.OpenId)
		if err != nil {
			return nil, err
		}
		point = pointInfo.Balance
	}

	return gin.H{
		"list":  res,
		"point": point,
		"total": total,
	}, nil
}

// CollectPrePoint  收集预加积分
func (receiver PlatformController) CollectPrePoint(c *gin.Context) (gin.H, error) {
	form := api.CollectPrePointRequest{}
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

	sceneUser := service.DefaultBdSceneUserService.FindPlatformUserByPlatformUserId(form.MemberId, form.PlatformKey)

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
	key := timeStr + ":prePoint:" + scene.Ch + sceneUser.PlatformUserId + userInfo.PhoneNumber

	lastPoint, _ := strconv.ParseInt(app.Redis.Get(ctx, key).Val(), 10, 64)
	incPoint := one.Point

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
	pointService := service.NewPointService(ctx)
	//积分
	types := entity.PointTypesMap[form.PlatformKey]
	point, err := pointService.IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:      sceneUser.OpenId,
		Type:        types,
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
		one.Point = halfPoint
	}

	err = repository.DefaultBdScenePrePointRepository.Save(&one)
	if err != nil {
		return nil, err
	}

	//减碳量
	fromString := decimal.NewFromInt(int64(one.Point))
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

func (receiver PlatformController) CheckMgs(c *gin.Context) (gin.H, error) {
	form := checkMsg{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)
	if form.Content != "" {
		//检查内容
		if err := validator.CheckMsgWithOpenId(user.OpenId, form.Content); err != nil {
			app.Logger.Errorf("文本校验 Error:%s\n", err.Error())
			zhuGeAttr := make(map[string]interface{}, 0)
			zhuGeAttr["场景"] = "文本校验"
			zhuGeAttr["失败原因"] = err.Error()
			track.DefaultZhuGeService().Track(config.ZhuGeEventName.MsgSecCheck, user.OpenId, zhuGeAttr)
			return nil, errno.ErrCommon.WithMessage(err.Error())
		}
	}
	return nil, nil
}

func (receiver PlatformController) CheckMedia(c *gin.Context) (gin.H, error) {
	form := checkMedia{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)
	if form.MediaUrl != "" {
		err := validator.CheckMediaWithOpenId(user.OpenId, form.MediaUrl)
		if err != nil {
			app.Logger.Errorf("图片校验 Error:%s\n", err.Error())
			zhuGeAttr := make(map[string]interface{}, 0)
			zhuGeAttr["场景"] = "图片校验"
			zhuGeAttr["失败原因"] = err.Error()
			track.DefaultZhuGeService().Track(config.ZhuGeEventName.MsgSecCheck, user.OpenId, zhuGeAttr)
			return nil, errno.ErrCommon.WithMessage(err.Error())
		}
	}
	return nil, nil
}

func (receiver PlatformController) bindCallback(scene entity.BdScene, sceneUser entity.BdSceneUser, user entity.User) error {
	//绑定回调
	var err error
	//var ch_key string
	if scene.Ch == "jinhuaxing" {
		params := make(map[string]interface{}, 0)
		params["mobile"] = user.PhoneNumber
		params["status"] = "1"
		jhxSvr := jhx.NewJhxService(context.NewMioContext())
		err = jhxSvr.BindSuccess(params)
		//ch_key = "金华行"
	}

	if scene.Ch == "yitongxing" {
		params := make(map[string]interface{}, 0)
		params["memberId"] = sceneUser.PlatformUserId
		params["openId"] = sceneUser.OpenId
		ytxSrv := ytx.NewYtxService(context.NewMioContext(), ytx.WithSecret(scene.Secret2), ytx.WithDomain(scene.Domain2))
		err = ytxSrv.BindSuccess(params)
		//ch_key = "亿通行"
	}

	return err
}
