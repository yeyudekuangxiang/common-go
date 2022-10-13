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
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	"mio/pkg/errno"
	"strconv"
	"time"
)

var DefaultJhxController = JhxController{}

type JhxController struct {
}

func (ctr JhxController) TicketCreate(ctx *gin.Context) (gin.H, error) {
	form := jhxTicketCreateRequest{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)
	jhxService := jhx.NewJhxService(context.NewMioContext())
	orderNo := "jhx" + strconv.FormatInt(time.Now().UnixNano(), 10)
	err := jhxService.TicketCreate(orderNo, 123, user)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

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

	params := make(map[string]string, 0)
	err := util.MapTo(&form, &params)
	if err != nil {
		return nil, err
	}
	sign := params["sign"]
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
func (ctr JhxController) JhxCollectPoint(ctx *gin.Context) (gin.H, error) {
	form := jhxCollectRequest{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	params := make(map[string]string, 0)
	err := util.MapTo(&form, &params)
	if err != nil {
		return nil, err
	}
	sign := params["sign"]
	delete(params, "sign")
	jhxService := jhx.NewJhxService(context.NewMioContext())
	point, err := jhxService.CollectPoint(sign, params)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"point": point,
	}, nil
}

func (ctr JhxController) JhxPreCollectPoint(c *gin.Context) (gin.H, error) {
	form := api.PreCollectRequest{}
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

	//查询用户
	userInfo, _ := service.DefaultUserService.GetUserBy(repository.GetUserBy{
		Mobile: form.Mobile,
		Source: entity.UserSourceMio,
	})

	//用户验证
	if userInfo.ID <= 0 {
		return nil, errno.ErrCommon.WithMessage("未找到用户")
	}

	//风险登记验证
	if userInfo.Risk >= 2 {
		fmt.Println("用户风险等级过高 ", form)
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
			err := jhx.NewJhxService(context.NewMioContext()).BindSuccess(sceneUser.Phone, "1")
			if err != nil {
				app.Logger.Errorf("callback jinhuaxing bind_success error:%s", err.Error())
			}
		}
	}

	//查重
	transService := service.NewPointTransactionService(ctx)
	typeString := service.DefaultBdSceneService.SceneToType(scene.Ch)

	by, err := transService.FindBy(repository.FindPointTransactionBy{
		OpenId: userInfo.OpenId,
		Type:   string(typeString),
		Note:   form.PlatformKey + "#" + form.Tradeno,
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
	amount, _ := fromString.Float64()
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
		})
		if err != nil {
			app.Logger.Errorf("预加积分 err:%s", err.Error())
		}
	}
	return gin.H{}, nil
}
