package open

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	"mio/pkg/errno"
	"strconv"
	"time"
)

var DefaultYtxController = YtxController{}

type YtxController struct {
}

func (ctr YtxController) AllReceive(ctx *gin.Context) (gin.H, error) {
	form := allReceiveRequest{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		app.Logger.Errorf("参数错误: %s", form)
		return nil, err
	}

	//查询 渠道信息
	scene := service.DefaultBdSceneService.FindByCh(form.PlatformKey)
	if scene.Key == "" || scene.Key == "e" {
		return nil, errno.ErrCommon.WithMessage("渠道查询失败")
	}

	//白名单验证
	ip := ctx.ClientIP()
	if err := service.DefaultBdSceneService.CheckWhiteList(ip, form.PlatformKey); err != nil {
		app.Logger.Info("校验白名单失败", ip)
		return nil, errno.ErrCommon.WithMessage("非白名单ip:" + ip)
	}

	user := apiutil.GetAuthUser(ctx)

	//风险登记验证
	if user.Risk >= 4 {
		return nil, errno.ErrCommon.WithMessage("账户风险等级过高")
	}

	sceneUser := service.DefaultBdSceneUserService.FindOne(repository.GetSceneUserOne{
		PlatformKey: form.PlatformKey,
		OpenId:      user.OpenId,
	})

	if sceneUser.ID == 0 {
		return nil, errno.ErrBindRecordNotFound
	}

	prePoint, total, err := repository.DefaultBdScenePrePointRepository.FindBy(repository.GetScenePrePoint{
		PlatformKey:    sceneUser.PlatformKey,
		PlatformUserId: sceneUser.PlatformUserId,
		//OpenId:         sceneUser.OpenId,
		Status: 1,
	})
	if err != nil {
		return nil, errno.ErrCommon.WithErr(err)
	}

	if total == 0 {
		return nil, errno.ErrRecordNotFound
	}

	//检查上限
	var halfPoint, incPoint, totalPoint, halfId, thanPoint int64
	var ids []int64

	timeStr := time.Now().Format("2006-01-02")
	key := timeStr + ":prePoint:" + scene.Ch + sceneUser.PlatformUserId + user.PhoneNumber

	totalPoint, _ = strconv.ParseInt(app.Redis.Get(ctx, key).Val(), 10, 64)

	if totalPoint >= int64(scene.PrePointLimit) {
		return nil, errno.ErrCommon.WithMessage("今日获取积分已达到上限")
	}

	for _, v := range prePoint {
		onePoint := v.Point
		totalPoint += onePoint

		if totalPoint > int64(scene.PrePointLimit) {
			thanPoint = totalPoint - int64(scene.PrePointLimit)

			halfPoint = onePoint - thanPoint
			halfId = v.ID

			totalPoint = int64(scene.PrePointLimit)
			break
		}
		incPoint += onePoint
		ids = append(ids, v.ID)
	}

	incPoint = incPoint + halfPoint

	if incPoint == 0 {
		return nil, nil
	}

	bizId := util.UUID()
	c := context.NewMioContext()
	_, err = service.NewPointService(c).IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       sceneUser.OpenId,
		Type:         entity.POINT_YTX,
		BizId:        bizId,
		ChangePoint:  incPoint,
		AdditionInfo: "一键领取亿通行积分",
	})

	if err != nil {
		return nil, errno.ErrCommon.WithErr(err)
	}

	app.Redis.Set(ctx, key, totalPoint, 24*time.Hour)

	typeCarbonStr := service.DefaultBdSceneService.SceneToCarbonType(form.PlatformKey)

	if typeCarbonStr != "" {
		_, err = service.NewCarbonTransactionService(c).Create(api_types.CreateCarbonTransactionDto{
			OpenId: user.OpenId,
			UserId: user.ID,
			Type:   typeCarbonStr,
			Value:  float64(incPoint / int64(scene.Override)),
			BizId:  bizId,
		})
		if err != nil {
			return nil, errno.ErrCommon.WithErr(err)
		}
	}

	upStatus := make(map[string]interface{}, 0)
	upStatus["status"] = 2
	if len(ids) > 0 {
		err = repository.DefaultBdScenePrePointRepository.Updates(repository.GetScenePrePoint{
			Ids:    ids,
			Status: 1,
		}, upStatus)
		if err != nil {
			return nil, errno.ErrCommon.WithErr(err)
		}

	}

	if halfId != 0 {
		upHalfPoint := make(map[string]interface{}, 0)
		upHalfPoint["point"] = thanPoint
		err = repository.DefaultBdScenePrePointRepository.Updates(repository.GetScenePrePoint{
			Id:     halfId,
			Status: 1,
		}, upHalfPoint)

		if err != nil {
			return nil, errno.ErrCommon.WithErr(err)
		}
	}

	return nil, nil
}

func (ctr YtxController) PrePointList(c *gin.Context) (gin.H, error) {
	form := PrePointListRequestV2{}
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

	user := apiutil.GetAuthUser(c)

	res, total, err := repository.DefaultBdScenePrePointRepository.FindBy(repository.GetScenePrePoint{
		PlatformKey: form.PlatformKey,
		OpenId:      user.OpenId,
		Status:      1,
		StartTime:   time.Now().AddDate(0, 0, -7),
	})

	if err != nil {
		return nil, errno.ErrCommon.WithErr(err)
	}

	pointInfo, err := service.NewPointService(context.NewMioContext()).FindByOpenId(user.OpenId)
	if err != nil {
		return nil, err
	}

	point := pointInfo.Balance

	return gin.H{
		"list":  res,
		"point": point,
		"total": total,
	}, nil
}
