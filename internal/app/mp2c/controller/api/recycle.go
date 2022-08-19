package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	"mio/pkg/errno"
	"strconv"
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
	if form.Type != 1 {
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

	if err := service.DefaultRecycleService.CheckSign(dst, scene.Key); err != nil {
		app.Logger.Info("校验sign失败", form)
		return nil, errors.New("sign:" + form.Sign + " 验证失败")
	}
	//避开重放
	if !util.DefaultLock.Lock(strconv.Itoa(form.Type)+form.OrderNo, 24*3600*30*time.Second) {
		fmt.Println("charge 重复提交订单", form)
		app.Logger.Info("charge 重复提交订单", form)
		return nil, errors.New("重复提交订单")
	}
	//通过openid查询用户
	userInfo, _ := service.DefaultUserService.GetUserByOpenId(form.ClientId)
	if userInfo.ID == 0 {
		fmt.Println("charge 未找到用户 ", form)
		return nil, errno.ErrUserNotFound
	}
	//查询今日获取积分次数
	limitKey := time.Now().Format("20060102") + form.Name + userInfo.OpenId
	if !util.DefaultLock.Lock(limitKey, time.Hour*24) {
		return nil, errors.New("今日该回收分类获取积分次数已达到上限")
	}
	//匹配类型
	typeName := service.DefaultRecycleService.GetType(form.Name)
	//查询今日积分上限
	currPoint := service.DefaultRecycleService.GetPoint(form.ProductCategoryName, form.Qua, form.Unit)
	maxPoint := service.DefaultRecycleService.GetMaxPointByMonth(form.ProductCategoryName)
	point, err := service.NewPointTransactionCountLimitService(context.NewMioContext()).
		CheckMaxPointByMonth(typeName, userInfo.OpenId, currPoint, maxPoint)
	if err != nil {
		return nil, err
	}
	//加积分
	pointService := service.NewPointService(context.NewMioContext())
	_, err = pointService.IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       userInfo.OpenId,
		Type:         typeName,
		ChangePoint:  point,
		BizId:        util.UUID(),
		AdditionInfo: form.OrderNo + "#" + form.ClientId,
	})
	if err != nil {
		fmt.Println("oola 旧物回收 加积分失败 ", form)
	}
	return gin.H{}, nil
}
