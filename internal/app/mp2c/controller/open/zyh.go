package open

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service/platform"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util/apiutil"
)

var DefaultZyhController = ZyhController{}

type ZyhController struct {
}

func (ctr ZyhController) SendPoint(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	zyhService := platform.NewZyhService(context.NewMioContext())
	println(user.ID)

	pointType := "1"
	point := "100"
	openid := "oy_BA5Nwkt6hzECxIXwNYkhLyzSs"

	messageCode, messageErr := zyhService.SendPoint(pointType, openid, point)

	zyhService.CreateLog(srv_types.GetZyhLogAddDTO{
		Openid:         openid,
		PointType:      entity.POINT_ARTICLE,
		Value:          1,
		ResultCode:     messageCode,
		AdditionalInfo: messageErr.Error(),
		TransactionId:  "1111",
	})

	if messageErr != nil {
		return nil, messageErr
	}
	//入库
	return nil, nil
}

//给客户用，查数据用

func (ctr ZyhController) Zyh(ctx *gin.Context) (gin.H, error) {
	form := api.GetZyhForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	if form.Mobile != "" {
		zyhService := platform.NewZyhService(context.NewMioContext())
		return zyhService.GetZyhInfoByMobile(srv_types.GetZyhOpenDTO{
			Mobile: form.Mobile,
		})
	}

	if form.VolId != "" {
		zyhService := platform.NewZyhService(context.NewMioContext())
		return zyhService.GetZyhInfoByVolId(srv_types.GetZyhOpenDTO{
			VolId: form.VolId,
		})
	}
	return gin.H{}, nil
}
