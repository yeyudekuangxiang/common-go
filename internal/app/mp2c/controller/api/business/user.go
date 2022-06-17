package business

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/service/business"
	"mio/internal/pkg/util/apiutil"
)

var DefaultUserController = UserController{}

type UserController struct{}

func (UserController) GetUserInfo(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthBusinessUser(ctx)
	//获取部门名称
	department, err := business.DefaultDepartmentService.GetBusinessDepartmentById(user.BDepartmentId)
	if err != nil {
		app.Logger.Error("部门信息查询失败", user)
	}
	totalCarbonReduce := business.DefaultCarbonCreditsLogService.GetUserTotalCarbonCreditsByUserId(user.ID)
	totalPoints := business.DefaultPointLogService.GetUserTotalPointsByUserId(user.ID)
	return gin.H{
		"info":              user,
		"department":        department,
		"totalCarbonReduce": totalCarbonReduce.Total,
		"totalPoints":       totalPoints.Total,
	}, nil
}

func (UserController) GetToken(ctx *gin.Context) (gin.H, error) {
	originInfo := apiutil.GetAuthUser(ctx)
	fmt.Println(originInfo)
	if originInfo.PhoneNumber == "" {
		return gin.H{
			"token": "",
		}, nil
	}
	//查询企业信息
	user, err := business.DefaultUserService.GetBusinessUserBy(business.GetBusinessUserParam{Mobile: originInfo.PhoneNumber})
	if err != nil {
		fmt.Println("GetToken 查询企业用户失败", err.Error())
		return gin.H{
			"token": "",
		}, nil
	}
	fmt.Println("user is ", user)
	//创建token
	token, err := business.DefaultUserService.CreateBusinessUserToken(user)
	if err != nil {
		fmt.Println("GetToken 创建token失败", err.Error())
		return gin.H{
			"token": "",
		}, nil
	}

	return gin.H{
		"token": token,
	}, nil
}
