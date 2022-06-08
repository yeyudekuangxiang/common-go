package business

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/service/business"
	"mio/internal/pkg/util/apiutil"
	"time"
)

var DefaultPointController = PointController{}

type PointController struct{}

func (PointController) GetPointRecordList(ctx *gin.Context) (gin.H, error) {
	form := GetPointRecordListForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthBusinessUser(ctx)

	infoList := business.DefaultPointLogService.GetPointLogInfoList(business.GetPointLogInfoListParam{
		UserId:    user.ID,
		StartTime: form.Date,
		EndTime:   form.Date.AddDate(0, 1, 0).Add(-time.Nanosecond),
	})
	return gin.H{
		"list": infoList,
	}, nil
}
