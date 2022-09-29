package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/app"
	mioctx "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/platform/jhx"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	"mio/pkg/errno"
	"strconv"
	"time"
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

//GetUserInfo 用户信息
func (UserController) GetUserInfo(c *gin.Context) (gin.H, error) {
	userInfo := apiutil.GetAuthUser(c)
	user := api_types.UserInfoVO{}
	util.MapTo(userInfo, &user)
	if userInfo.ChannelId != 0 {
		channel, err := service.DefaultUserChannelService.GetChannelInfoByCid(userInfo.ChannelId)
		if err == nil {
			user.ChannelName = channel.Name
		}
	}
	return gin.H{
		"user": user,
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
		user, err := service.DefaultUserService.FindOrCreateByMobile(form.Mobile, form.Cid)
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

func (UserController) BindMobileByYZM(c *gin.Context) (gin.H, error) {
	form := GetYZMForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	userInfo := apiutil.GetAuthUser(c)
	if service.DefaultUserService.CheckYZM(form.Mobile, form.Code) {
		return nil, service.DefaultUserService.BindMobileByYZM(userInfo.ID, form.Mobile)
	}
	return nil, errors.New("验证码错误,请重新输入")
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
	//绑定后
	err := service.DefaultUserService.BindPhoneByCode(user.ID, form.Code, c.ClientIP(), form.InvitedBy)
	if err == nil && user.ChannelId == 1059 {
		go func() {
			jhxService := jhx.NewJhxService(mioctx.NewMioContext())
			orderNo := "jhx" + strconv.FormatInt(time.Now().Unix(), 10)
			startTime, _ := time.Parse("2006-01-02", "2022-09-29")
			endTime, _ := time.Parse("2006-01-02", "2022-10-31")
			for i := 0; i < 2; i++ {
				err := jhxService.TicketCreate(orderNo+strconv.Itoa(i), 123, startTime, endTime, user)
				if err != nil {
					app.Logger.Errorf("金华行发券失败:%s", err.Error())
					return
				}
			}
			return
		}()
	}
	return nil, service.DefaultUserService.BindPhoneByCode(user.ID, form.Code, c.ClientIP(), form.InvitedBy)
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
func (UserController) GetUserAccountInfo(c *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(c)
	accountInfo, err := service.DefaultUserService.AccountInfo(user.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"accountInfo": accountInfo,
	}, nil
}

func (UserController) UpdateUserInfo(c *gin.Context) (gin.H, error) {
	form := UpdateUserInfoForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)

	var birthday *time.Time
	if form.Birthday != nil {
		t, err := time.Parse("2006-01-02", *form.Birthday)
		if err != nil {
			return nil, errno.ErrBind.WithErr(err)
		}
		birthday = &t
	}

	var gender *entity.UserGender
	if form.Gender != nil {
		gender = (*entity.UserGender)(form.Gender)
	}

	err := service.DefaultUserService.UpdateUserInfo(service.UpdateUserInfoParam{
		UserId:      user.ID,
		Nickname:    form.Nickname,
		Avatar:      form.Avatar,
		Gender:      gender,
		Birthday:    birthday,
		PhoneNumber: form.PhoneNumber,
	})
	return nil, err
}
