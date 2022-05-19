package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
)

var DefaultUserController = UserController{}

type UserController struct {
}

func (UserController) GetNewUser(c *gin.Context) (gin.H, error) {
	user, err := service.DefaultUserService.GetUserById(1)
	f, err := excelize.OpenFile("/Users/leo/Downloads/test1.xlsx")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Get value from cell by given worksheet name and axis.
	cell, err := f.GetCellValue("sheet1", "A2")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(cell)
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("sheet1")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell)
		}
		fmt.Println()
	}

	if err != nil {
		return nil, err
	}
	return gin.H{
		"user": user,
	}, nil
}

func (UserController) GetUserInfo(c *gin.Context) (gin.H, error) {
	return gin.H{
		"user": apiutil.GetAuthUser(c),
	}, nil
}

func (UserController) GetYZM(c *gin.Context) (gin.H, error) {
	form := GetYZMForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	_, err := service.DefaultUserService.GetYZM(form.Mobile)
	if err != nil {
		return gin.H{
			"msg": "fail",
		}, err
	}
	return nil, nil
}

func (UserController) CheckYZM(c *gin.Context) (gin.H, error) {
	form := GetYZMForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	if service.DefaultUserService.CheckYZM(form.Mobile, form.Code) {
		user, err := service.DefaultUserService.FindOrCreateByMobile(form.Mobile)
		if err != nil {
			return gin.H{}, err
		}
		userId := user.ID
		token, err := service.DefaultUserService.CreateUserToken(userId)
		return gin.H{
			"token":  token,
			"userId": userId,
		}, err
	} else {
		err := errors.New("验证码错误,请重新输入")
		return gin.H{}, err
	}

}

func (UserController) GetMobileUserInfo(c *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(c)
	mobileUser, err := service.DefaultUserService.FindUserBySource(entity.UserSourceMobile, user.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"user": mobileUser.ShortUser(),
	}, nil
}

func (UserController) BindMobileByCode(c *gin.Context) (gin.H, error) {
	form := BindMobileByCodeForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)
	return nil, service.DefaultUserService.BindPhoneByCode(user.ID, form.Code)
}
func (UserController) GetUserSummary(c *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(c)
	summary, err := service.DefaultUserService.UserSummary(user.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"summary": summary,
	}, nil
}

func (UserController) UpdateUserInfo(c *gin.Context) (gin.H, error) {
	form := UpdateUserInfoForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)

	var gender entity.UserGender
	if form.Gender == 1 {
		gender = entity.UserGenderMale
	} else if form.Gender == 2 {
		gender = entity.UserGenderFemale
	}
	err := service.DefaultUserService.UpdateUserInfo(service.UpdateUserInfoParam{
		UserId:   user.ID,
		Nickname: form.Nickname,
		Avatar:   form.Avatar,
		Gender:   gender,
	})
	return nil, err
}
