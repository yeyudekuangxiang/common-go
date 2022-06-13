package business

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/service/business"
	"mio/internal/pkg/util/apiutil"
)

var DefaultUserController = UserController{}

type UserController struct{}

func (UserController) GetUserInfo(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthBusinessUser(ctx)
	//先拿token,然后通过手机号和公司id差信息
	return gin.H{
		"info": user,
	}, nil
}

func (UserController) GetToken(ctx *gin.Context) (gin.H, error) {
	originInfo := apiutil.GetAuthUser(ctx)
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
	//创建token
	token, _ := business.DefaultUserService.CreateBusinessUserToken(user)
	return gin.H{
		"token": token,
	}, nil
}
