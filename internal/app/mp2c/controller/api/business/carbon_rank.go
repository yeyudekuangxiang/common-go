package business

import (
	"github.com/gin-gonic/gin"
	ebusiness "mio/internal/pkg/model/entity/business"
	"mio/internal/pkg/service/business"
	"mio/internal/pkg/util/apiutil"
)

var DefaultCarbonRankController = CarbonRankController{}

type CarbonRankController struct{}

func (CarbonRankController) GetUserRankList(ctx *gin.Context) (gin.H, error) {
	form := GetUserRankListForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthBusinessUser(ctx)

	list, total, err := business.DefaultCarbonRankService.UserRankList(business.GetUserRankListParam{
		UserId:    user.ID,
		DateType:  ebusiness.RankDateType(form.DateType),
		CompanyId: user.BCompanyId,
		Limit:     form.Limit(),
		Offset:    form.Offset(),
	})
	if err != nil {
		return nil, err
	}

	return gin.H{
		"list":     list,
		"total":    total,
		"page":     form.Page,
		"pageSize": form.PageSize,
	}, nil
}
