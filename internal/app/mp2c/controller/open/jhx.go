package open

import (
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
	platformUtil "mio/internal/pkg/util/platform"
	"mio/pkg/errno"
	"strconv"
	"time"
)

var DefaultJhxController = JhxController{}

type JhxController struct {
}

//发放券码
func (ctr JhxController) TicketCreate(ctx *gin.Context) (gin.H, error) {
	form := jhxTicketCreateRequest{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)
	jhxService := jhx.NewJhxService(context.NewMioContext())
	tradeNo, err := jhxService.SendCoupon(1000, user)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"orderNo": tradeNo,
	}, nil
}

//查询券码状态
func (ctr JhxController) TicketStatus(ctx *gin.Context) (gin.H, error) {
	form := jhxTicketStatusRequest{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	jhxService := jhx.NewJhxService(context.NewMioContext())
	result, err := jhxService.TicketStatus(form.Tradeno)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"status":   result.Status,
		"usedTime": result.UsedTime,
	}, nil
}

//消费通知
func (ctr JhxController) BusTicketNotify(ctx *gin.Context) (gin.H, error) {
	form := jhxTicketNotifyRequest{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	params := make(map[string]interface{}, 0)
	err := util.MapTo(&form, &params)
	if err != nil {
		return nil, err
	}
	sign := form.Sign
	delete(params, "sign")

	jhxService := jhx.NewJhxService(context.NewMioContext())
	err = jhxService.TicketNotify(sign, params)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//获取积分气泡list
func (ctr JhxController) JhxGetPreCollectPoint(ctx *gin.Context) (gin.H, error) {
	form := jhxGetPreCollectRequest{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	//查询 渠道信息
	scene := service.DefaultBdSceneService.FindByCh(form.PlatformKey)
	if scene.Key == "" || scene.Key == "e" {
		app.Logger.Info("渠道查询失败", form)
		return nil, errno.ErrChannelNotFound
	}

	params := make(map[string]interface{}, 0)
	err := util.MapTo(&form, &params)
	if err != nil {
		return nil, err
	}
	sign := form.Sign
	delete(params, "sign")

	jhxService := jhx.NewJhxService(context.NewMioContext())
	item, point, err := jhxService.GetPreCollectPointList(sign, params)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"list":  item,
		"point": point,
	}, nil
}

//消费积分气泡
func (ctr JhxController) JhxCollectPoint(c *gin.Context) (gin.H, error) {
	form := jhxCollectRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	scene := repository.DefaultBdSceneRepository.FindByCh(form.PlatformKey)
	if scene.Key == "" || scene.Key == "e" {
		return nil, errno.ErrCommon.WithMessage("渠道查询失败")
	}

	sceneUser := repository.DefaultBdSceneUserRepository.FindOne(repository.GetSceneUserOne{
		PlatformKey:    form.PlatformKey,
		PlatformUserId: form.MemberId,
		OpenId:         form.OpenId,
	})
	if sceneUser.ID == 0 {
		return nil, errno.ErrCommon.WithMessage("未找到绑定关系")
	}

	params := make(map[string]interface{}, 0)
	err := util.MapTo(&form, &params)
	if err != nil {
		return nil, err
	}

	sign := form.Sign
	delete(params, "sign")

	if err = platformUtil.CheckSign(sign, params, "", "&"); err != nil {
		return nil, err
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))

	var collect jhx.Collect
	_ = util.MapTo(&form, &collect)

	resp, err := jhx.NewJhxService(ctx).CollectPoint(collect)

	if err != nil {
		return nil, err
	}

	//检查上限
	var isHalf bool
	var halfPoint int64
	timeStr := time.Now().Format("2006-01-02")
	key := timeStr + ":prePoint:" + form.PlatformKey + form.MemberId
	lastPoint, _ := strconv.ParseInt(app.Redis.Get(ctx.Context, key).Val(), 10, 64)

	incPoint := resp.Point

	totalPoint := lastPoint + incPoint

	limit := int64(scene.PrePointLimit)
	if lastPoint >= limit {
		return nil, errno.ErrCommon.WithMessage("今日获取积分已达到上限")
	}

	if totalPoint > limit {
		p := incPoint
		incPoint = limit - lastPoint
		totalPoint = limit
		isHalf = true
		halfPoint = p - incPoint
	}

	app.Redis.Set(ctx.Context, key, totalPoint, 24*time.Hour)

	//积分
	point, err := service.NewPointService(context.NewMioContext()).IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:      sceneUser.OpenId,
		Type:        entity.POINT_JHX,
		BizId:       util.UUID(),
		ChangePoint: incPoint,
		AdminId:     0,
		Note:        form.PlatformKey + "#" + resp.Tradeno,
	})

	if err != nil {
		return nil, err
	}

	//更新pre_point对应数据
	resp.Status = 2
	resp.UpdatedAt = time.Now()
	if isHalf {
		resp.Status = 1
		resp.Point = halfPoint
	}

	err = service.NewBdScenePrePointService(ctx).Save(resp)

	if err != nil {
		return nil, err
	}

	return gin.H{
		"point": point,
	}, nil
}

//金华行单独调用 预加积分
func (ctr JhxController) JhxPreCollectPoint(c *gin.Context) (gin.H, error) {
	form := api.PreCollectRequest{}
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
		return nil, errno.ErrCommon.WithMessage("sign:" + form.Sign + " 验证失败")
	}

	//查询用户
	userInfo, exist, err := service.DefaultUserService.GetUser(repository.GetUserBy{
		Mobile: form.Mobile,
		Source: entity.UserSourceMio,
	})
	if err != nil {
		return nil, errno.ErrCommon.WithMessage(err.Error())
	}
	//用户验证
	if !exist {
		return nil, errno.ErrUserNotFound
	}

	//查重
	_, exist, err = repository.DefaultBdScenePrePointRepository.FindOne(repository.GetScenePrePoint{
		PlatformKey: form.PlatformKey,
		TradeNo:     form.Tradeno,
	})
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, errno.ErrCommon.WithMessage("重复提交订单")
	}

	//入参保存
	defer trackBehaviorInteraction(trackInteractionParam{
		Tp:     form.PlatformKey,
		Data:   form,
		Ip:     c.ClientIP(),
		UserId: userInfo.ID,
	})

	//风险登记验证
	if userInfo.Risk >= 4 {
		return nil, errno.ErrCommon.WithMessage("账户风险等级过高")
	}

	sceneUser := repository.DefaultBdSceneUserRepository.FindOne(repository.GetSceneUserOne{
		PlatformKey: form.PlatformKey,
		//PlatformUserId: form.MemberId,
		OpenId: userInfo.OpenId,
	})

	//第三方推送订单时检测是否绑定，未绑定用户执行绑定
	if sceneUser.ID == 0 {
		sceneUser.PlatformKey = scene.Ch
		sceneUser.PlatformUserId = form.MemberId
		sceneUser.Phone = userInfo.PhoneNumber
		sceneUser.OpenId = userInfo.OpenId
		sceneUser.UnionId = userInfo.UnionId
		err = service.DefaultBdSceneUserService.Create(sceneUser)
		if err != nil {
			app.Logger.Errorf("bind db_scene_user error:%s", err.Error())
		}
		//绑定回调
		if form.PlatformKey == "jinhuaxing" {
			bindParams := make(map[string]interface{}, 0)
			bindParams["mobile"] = userInfo.PhoneNumber
			bindParams["status"] = "1"
			err := jhx.NewJhxService(context.NewMioContext()).BindSuccess(bindParams)
			if err != nil {
				app.Logger.Errorf("callback jinhuaxing bind_success error:%s", err.Error())
			}
		}
	}

	//预加积分
	fromString, _ := decimal.NewFromString(form.Amount)
	amount, _ := fromString.Float64()

	point := fromString.Mul(decimal.NewFromInt(int64(scene.Override))).Round(0).IntPart()
	err = repository.DefaultBdScenePrePointRepository.Create(&entity.BdScenePrePoint{
		PlatformKey:    form.PlatformKey,
		PlatformUserId: form.MemberId,
		Point:          point,
		OpenId:         userInfo.OpenId,
		Status:         1,
		Mobile:         userInfo.PhoneNumber,
		Tradeno:        form.Tradeno,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		return nil, err
	}
	//减碳量
	typeCarbonStr := service.DefaultBdSceneService.SceneToCarbonType(scene.Ch)
	if typeCarbonStr != "" {
		_, err = service.NewCarbonTransactionService(context.NewMioContext()).Create(api_types.CreateCarbonTransactionDto{
			OpenId: userInfo.OpenId,
			UserId: userInfo.ID,
			Type:   typeCarbonStr,
			Value:  amount,
			Ip:     ip,
			BizId:  form.Tradeno,
		})
		if err != nil {
			app.Logger.Errorf("预加积分 err:%s", err.Error())
		}
	}
	return gin.H{}, nil
}
