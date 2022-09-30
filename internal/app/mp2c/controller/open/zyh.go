package open

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/model/entity/activity"
	"mio/internal/pkg/service/platform"
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

	zyhService.SendPointLog(activity.ZyhLog{
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
