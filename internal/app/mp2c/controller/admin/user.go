package admin

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
)

var DefaultUserController = UserController{}

type UserController struct {
}

func (UserController) GetUserInfo(c *gin.Context) (gin.H, error) {
	var form GetUserForm
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user, err := service.DefaultUserService.GetUserById(form.Id)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"user": user,
	}, nil
}

func GetUserPageListBy(c *gin.Context) (gin.H, error) {
	form := UserPageListForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	list, count := service.DefaultUserService.GetUserPageListBy(repository.GetUserPageListBy{
		Limit:   10,
		Offset:  0,
		User:    repository.GetUserListBy{Mobile: form.Mobile},
		OrderBy: "id desc",
	})
	return gin.H{
		"users": list,
		"page":  count,
	}, nil
}
