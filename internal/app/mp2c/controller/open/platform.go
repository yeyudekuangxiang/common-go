package open

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"mio/internal/app/mp2c/controller/api"
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
		return nil, errors.New("用户不存在")
	}

	method := scene.Ch
	if form.Method != "" {
		method = strings.ToLower(method) + "_" + strings.ToLower(form.Method)
	}
	if _, ok := entity.PlatformMethodMap[method]; !ok {
		return nil, errors.New("未找到匹配方法")
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
		return nil, errors.New("渠道查询失败")
	}

	//白名单验证
	ip := c.ClientIP()
	if err := service.DefaultBdSceneService.CheckWhiteList(ip, form.PlatformKey); err != nil {
		app.Logger.Info("校验白名单失败", ip)
		return nil, errors.New("非白名单ip:" + ip)
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
		return nil, errors.New("sign:" + form.Sign + " 验证失败")
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

	//预加积分
	fromString, _ := decimal.NewFromString(params["amount"])
	point := fromString.Mul(decimal.NewFromInt(int64(scene.Override))).Round(0).String()
	err = repository.DefaultBdScenePrePointRepository.Create(&entity.BdScenePrePoint{
		PlatformKey:    form.PlatformKey,
		PlatformUserId: form.MemberId,
		Point:          point,
		Status:         1,
		Tradeno:        form.TradeNo,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})

	if err != nil {
		return nil, errno.ErrCommon.WithErr(err)
	}

	return gin.H{}, nil
}

// GetPrePointList 获取预加积分列表
func (receiver PlatformController) GetPrePointList(c *gin.Context) (gin.H, error) {

	return nil, nil
}

// CollectPoint 收集预加积分
func (receiver PlatformController) CollectPoint(c *gin.Context) (gin.H, error) {

	////减碳量
	//fromString, _ := decimal.NewFromString(params["amount"])
	//amount, _ := fromString.Float64()
	//typeCarbonStr := service.DefaultBdSceneService.SceneToCarbonType(scene.Ch)
	//if typeCarbonStr != "" {
	//	_, err = service.NewCarbonTransactionService(context.NewMioContext()).Create(api_types.CreateCarbonTransactionDto{
	//		OpenId: sceneUser.OpenId,
	//		UserId: userInfo.ID,
	//		Type:   typeCarbonStr,
	//		Value:  amount,
	//		Ip:     ip,
	//	})
	//	if err != nil {
	//		app.Logger.Errorf("预加积分 err:%s", err.Error())
	//	}
	//}
	return nil, nil
}
