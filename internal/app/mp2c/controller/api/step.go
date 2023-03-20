package api

import (
	"github.com/gin-gonic/gin"
	"github.com/medivhzhan/weapp/v3"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
	"sort"
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
func (StepController) UpdateStep(ctx *gin.Context) (gin.H, error) {
	form := UpdateStepForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthUser(ctx)

	stepList := make([]weapp.SetpInfo, 0)
	for _, item := range form.StepList {
		stepList = append(stepList, weapp.SetpInfo{
			Step:      item.Step,
			Timestamp: item.Timestamp,
		})
	}
	sort.Slice(stepList, func(i, j int) bool {
		return stepList[i].Timestamp < stepList[j].Timestamp
	})
	err := service.DefaultStepHistoryService.UpdateStepHistoryByList(user.OpenId, stepList)
	return nil, err
}
func (StepController) WeeklyHistory(ctx *gin.Context) (interface{}, error) {
	user := apiutil.GetAuthUser(ctx)
	return service.DefaultStepService.WeeklyHistory(user.OpenId)
}
func (StepController) Collect(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	carbon, err := service.DefaultStepService.RedeemCarbonFromPendingSteps(user.OpenId, user.ID, ctx.ClientIP())
	if err != nil {
		return gin.H{
			"points": 0,
			"carbon": carbon,
		}, err
	}
	var points int
	points, err = service.DefaultStepService.RedeemPointFromPendingSteps(user.OpenId)
	return gin.H{
		"points": points,
		"carbon": carbon,
	}, err
}
