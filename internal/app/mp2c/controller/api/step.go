package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
)

var DefaultStepController = StepController{}

type StepController struct {
}

func (StepController) UpdateStepTotal(ctx *gin.Context) (gin.H, error) {
	form := UpdateStepTotalForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)

	err := service.DefaultStepHistoryService.UpdateStepHistoryByEncrypted(service.UpdateStepHistoryByEncryptedParam{
		OpenId:        user.OpenId,
		EncryptedData: form.EncryptedData,
		IV:            form.IV,
	})
	return nil, err
}

func (StepController) WeeklyHistory(ctx *gin.Context) (interface{}, error) {
	user := apiutil.GetAuthUser(ctx)
	return service.DefaultStepService.WeeklyHistory(user.OpenId)
}
func (StepController) Collect(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	points, err := service.DefaultStepService.RedeemPointFromPendingSteps(user.OpenId)
	return gin.H{
		"points": points,
	}, err
}
