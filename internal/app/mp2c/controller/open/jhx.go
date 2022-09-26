package open

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/platform"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	"strconv"
	"time"
)

var DefaultJhxController = JhxController{}

type JhxController struct {
}

func (ctr JhxController) TicketCreate(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	jhxService := platform.NewJhxService(context.NewMioContext())
	orderNo := "jhx" + strconv.FormatInt(time.Now().UnixNano(), 10)
	err := jhxService.TicketCreate(orderNo, user)
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
	jhxService := platform.NewJhxService(context.NewMioContext())
	result, err := jhxService.TicketStatus(form.Tradeno)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"status":   result.Status,
		"usedTime": result.UsedTime,
	}, nil
}

func (ctr JhxController) BusTicketNotify(ctx *gin.Context) (gin.H, error) {
	form := jhxTicketNotifyRequest{}
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
	jhxService := platform.NewJhxService(context.NewMioContext())
	err = jhxService.TicketNotify(sign, params)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//生产积分气泡
func (ctr JhxController) PreCollectPoint(ctx *gin.Context) (gin.H, error) {
	form := jhxCollectRequest{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	scene := service.DefaultBdSceneService.FindByCh("jhx")

	params := make(map[string]string, 0)
	err := util.MapTo(&form, &params)
	if err != nil {
		return nil, err
	}
	params["key"] = scene.Key
	sign := params["sign"]
	delete(params, "sign")

	jhxService := platform.NewJhxService(context.NewMioContext())
	err = jhxService.PreCollectPoint(sign, params)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//获取积分气泡
func (ctr JhxController) GetPreCollectPoint(ctx *gin.Context) (gin.H, error) {
	form := jhxGetCollectRequest{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	scene := service.DefaultBdSceneService.FindByCh("jhx")

	params := make(map[string]string, 0)
	err := util.MapTo(&form, &params)
	if err != nil {
		return nil, err
	}
	params["key"] = scene.Key
	sign := params["sign"]
	delete(params, "sign")

	jhxService := platform.NewJhxService(context.NewMioContext())
	list, err := jhxService.GetPreCollectPointList(sign, params)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"list": list,
	}, nil
}

//获取积分气泡
func (ctr JhxController) CollectPoint(ctx *gin.Context) (gin.H, error) {
	form := jhxCollectRequest{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	scene := service.DefaultBdSceneService.FindByCh("jhx")

	params := make(map[string]string, 0)
	err := util.MapTo(&form, &params)
	if err != nil {
		return nil, err
	}
	params["key"] = scene.Key
	sign := params["sign"]
	delete(params, "sign")
	jhxService := platform.NewJhxService(context.NewMioContext())
	err = jhxService.CollectPoint(sign, params)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
