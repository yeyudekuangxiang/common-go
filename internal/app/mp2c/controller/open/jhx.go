package open

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/context"
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
	//入库

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
	params := make(map[string]interface{}, 0)
	sign := params["sign"].(string)
	delete(params, "sign")
	err := util.MapTo(&form, &params)
	if err != nil {
		return nil, err
	}
	jhxService := platform.NewJhxService(context.NewMioContext())
	jhxService.TicketNotify(sign, params)
	return nil, nil
}
