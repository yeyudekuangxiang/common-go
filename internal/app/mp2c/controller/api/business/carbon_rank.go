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

	myRank, err := business.DefaultCarbonRankService.FindUserRank(business.FindUserRankParam{
		UserId:   user.ID,
		DateType: ebusiness.RankDateType(form.DateType),
	})
	if err != nil {
		return nil, err
	}

	return gin.H{
		"list":     list,
		"myRank":   myRank,
		"total":    total,
		"page":     form.Page,
		"pageSize": form.PageSize,
	}, nil
}
func (CarbonRankController) GetDepartmentRankList(ctx *gin.Context) (gin.H, error) {
	form := GetDepartmentRankListForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthBusinessUser(ctx)

	list, total, err := business.DefaultCarbonRankService.DepartmentRankList(business.GetDepartmentRankListParam{
		UserId:    user.ID,
		DateType:  ebusiness.RankDateType(form.DateType),
		CompanyId: user.BCompanyId,
		Limit:     form.Limit(),
		Offset:    form.Offset(),
	})
	if err != nil {
		return nil, err
	}

	myDepartmentRank, err := business.DefaultCarbonRankService.FindDepartmentRank(business.FindDepartmentRankParam{
		UserId:       user.ID,
		DepartmentId: user.BDepartmentId,
		DateType:     ebusiness.RankDateType(form.DateType),
	})
	if err != nil {
		return nil, err
	}

	return gin.H{
		"list":             list,
		"myDepartmentRank": myDepartmentRank,
		"total":            total,
		"page":             form.Page,
		"pageSize":         form.PageSize,
	}, nil
}
func (CarbonRankController) ChangeUserLikeStatus(ctx *gin.Context) (gin.H, error) {
	form := ChangeUserLikeStatusForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthBusinessUser(ctx)
	like, err := business.DefaultCarbonRankService.ChangeLikeStatus(business.ChangeLikeStatusParam{
		Pid:        form.Pid,
		ObjectType: form.ObjectType,
		DateType:   form.DateType,
		UserId:     user.ID,
	})
	if err != nil {
		return nil, err
	}
	return gin.H{
		"likeStatus": like.Status.IsLike(),
	}, nil
}
