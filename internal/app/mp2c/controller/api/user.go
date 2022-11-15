package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/app"
	mioctx "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/common"
	"mio/internal/pkg/service/platform/jhx"
	"mio/internal/pkg/service/platform/ytx"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	"mio/pkg/errno"
	"strings"
	"time"
)

var DefaultUserController = UserController{}

type UserController struct {
}

func (ctr UserController) GetNewUser(c *gin.Context) (gin.H, error) {
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
func (ctr UserController) GetUserInfo(c *gin.Context) (gin.H, error) {
	userInfo := apiutil.GetAuthUser(c)
	user := api_types.UserInfoVO{}
	_ = util.MapTo(userInfo, &user)
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

func (ctr UserController) GetYZM(c *gin.Context) (gin.H, error) {
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

func (ctr UserController) CheckYZM(c *gin.Context) (gin.H, error) {
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
		err := errno.ErrCommon.WithMessage("验证码错误,请重新输入")
		return gin.H{}, err
	}
}

func (ctr UserController) CheckYZMCommon(c *gin.Context) (gin.H, error) {
	form := GetYZMForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	if service.DefaultUserService.CheckYZM(form.Mobile, form.Code) {
		return gin.H{}, nil
	} else {
		err := errno.ErrCommon.WithMessage("验证码错误,请重新输入")
		return gin.H{}, err
	}
}

func (ctr UserController) BindMobileByYZM(c *gin.Context) (gin.H, error) {
	form := GetYZMForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	userInfo := apiutil.GetAuthUser(c)
	if service.DefaultUserService.CheckYZM(form.Mobile, form.Code) {
		return nil, service.DefaultUserService.BindMobileByYZM(userInfo.ID, form.Mobile)
	}
	return nil, errno.ErrCommon.WithMessage("验证码错误,请重新输入")
}

func (ctr UserController) GetMobileUserInfo(c *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(c)
	mobileUser, err := service.DefaultUserService.FindUserBySource(entity.UserSourceMobile, user.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"user": mobileUser.ShortUser(),
	}, nil
}

func (ctr UserController) BindMobileByCode(c *gin.Context) (gin.H, error) {
	form := BindMobileByCodeForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)
	//绑定后
	err := service.DefaultUserService.BindPhoneByCode(user.ID, form.Code, c.ClientIP(), form.InvitedBy)
	if err == nil {
		ctr.sendCoupon(user)
	}
	return nil, err
}

func (ctr UserController) GetUserSummary(c *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(c)
	summary, err := service.DefaultUserService.UserSummary(user.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"summary": summary,
	}, nil
}

func (ctr UserController) GetUserAccountInfo(c *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(c)
	accountInfo, err := service.DefaultUserService.AccountInfo(user.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"accountInfo": accountInfo,
	}, nil
}

func (ctr UserController) UpdateUserInfo(c *gin.Context) (gin.H, error) {
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

// 个人主页
func (ctr UserController) HomePage(c *gin.Context) (gin.H, error) {
	// 头像 昵称 笔记数量 ip属地 简介
	form := HomePageRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return gin.H{}, err
	}
	user := apiutil.GetAuthUser(c)

	if form.UserId != 0 {
		u, b, err := service.DefaultUserService.GetUserByID(form.UserId)
		if err != nil {
			return gin.H{}, err
		}

		if !b {
			return gin.H{}, nil
		}
		user = *u
	}

	//归属地
	location, err := common.NewCityService(mioctx.NewMioContext()).GetByCityCode(common.GetByCityCodeParams{CityCode: user.CityCode})
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	shortUser := user.ShortUser()
	result := make(map[string]interface{}, 0)
	_ = util.MapTo(&shortUser, &result)
	result["ipLocation"] = location.Name
	return result, nil
}

// 更新简介
func (ctr UserController) UpdateIntroduction(c *gin.Context) (gin.H, error) {
	form := UpdateIntroductionRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return gin.H{}, err
	}

	user := apiutil.GetAuthUser(c)

	if len(strings.Trim(form.Introduction, " ")) == 0 {
		return nil, errno.ErrCommon.WithMessage("内容不可以为空")
	}

	err := service.DefaultUserService.UpdateUserInfo(service.UpdateUserInfoParam{
		UserId:       user.ID,
		Introduction: form.Introduction,
	})

	return nil, err
}

func (ctr UserController) sendCoupon(user entity.User) {
	if user.ChannelId == 1059 {
		jhxService := jhx.NewJhxService(mioctx.NewMioContext())
		go func() {
			for i := 0; i < 2; i++ {
				_, err := jhxService.SendCoupon(1000, user)
				if err != nil {
					app.Logger.Errorf("金华行发券失败:%s", err.Error())
					return
				}
			}
			return
		}()
	}

	if user.ChannelId == 1066 {
		bdscene := service.DefaultBdSceneService.FindByCh("yitongxing")
		var options []ytx.Options
		options = append(options, ytx.WithPoolCode(bdscene.AppId2))
		options = append(options, ytx.WithSecret(bdscene.Secret))
		options = append(options, ytx.WithAppId(bdscene.AppId))
		options = append(options, ytx.WithDomain(bdscene.Domain))
		ytxService := ytx.NewYtxService(mioctx.NewMioContext(), options...)
		go func() {
			app.Logger.Info("亿通行发红包开始")
			_, err := ytxService.SendCoupon(1001, 5.00, user)
			if err != nil {
				app.Logger.Errorf("亿通行发红包失败:%s", err.Error())
				return
			}
			app.Logger.Info("亿通行发红包结束")
		}()
	}
}
