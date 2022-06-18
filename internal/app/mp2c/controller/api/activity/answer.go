package activity

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/activity"
	"mio/internal/pkg/util/apiutil"
	"strconv"
	"time"
)

var DefaultAnswerController = AnswerController{}

type AnswerController struct {
}

// HomePage Answer the questions
func (ctr AnswerController) HomePage(ctx *gin.Context) (gin.H, error) {
	//当前登陆用户信息
	user := apiutil.GetAuthUser(ctx)
	//if user.ID != 0 && time.Now().Add(-24*time.Hour).After(user.Time.Time) {
	//	return nil, errors.New("本次活动仅限绿喵新用户参加哦～")
	//}

	form := &GDDbActivityHomePageForm{}
	if err := apiutil.BindForm(ctx, form); err != nil {
		return nil, err
	}

	res, err := activity.DefaultGDdbService.HomePage(user.ID, form.InviteId)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"record": res,
	}, err
}

// GetUserSchool 获取用户及学校信息
func (ctr AnswerController) GetUserSchool(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	school, err := activity.DefaultGDdbService.GetUserSchool(user.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"record": school,
	}, nil
}

func (ctr AnswerController) PutFile(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	t, _ := strconv.Atoi(ctx.PostForm("type"))
	file, err := ctx.FormFile("file")
	if err != nil {
		return nil, err
	}
	if file.Size > 5*1024*1024 {
		return nil, errors.New("上传失败，请调整文件大小！")
	}
	open, err := file.Open()
	if err != nil {
		return nil, err
	}
	object, err := service.DefaultOssService.PutObject(fmt.Sprintf("/activity/gd/%d_%d", user.ID, time.Now().Unix()), open)
	if err != nil {
		return nil, err
	}
	// 更新
	err = activity.DefaultGDdbService.UpdateActivityUser(user.ID, t, object)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"object": object,
	}, nil
}

//StartQuestion 开始答题
func (ctr AnswerController) StartQuestion(ctx *gin.Context) (gin.H, error) {
	// 更新答题状态
	user := apiutil.GetAuthUser(ctx)
	err := activity.DefaultGDdbService.UpdateAnswerStatus(user.ID, 1)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return nil, nil
}

//EndQuestion 答题完成
func (ctr AnswerController) EndQuestion(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	form := &GDDbActivitySchoolForm{}
	if err := apiutil.BindForm(ctx, form); err != nil {
		return nil, err
	}

	// 保存学校信息，更新答题状态
	err := activity.DefaultGDdbService.SaveSchoolInfo(form.UserName, form.SchoolId, form.GradeId, user.ID, form.ClassNumber)
	if err != nil {
		return nil, errors.New("学校信息保存失败")
	}
	// 检测成团状态
	err = activity.DefaultGDdbService.CheckActivityStatus(user.ID, form.SchoolId)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// GetCityList 获取城市列表
func (ctr AnswerController) GetCityList(ctx *gin.Context) (gin.H, error) {
	cityList := activity.DefaultGDdbService.GetCityList()
	return gin.H{
		"record": cityList,
	}, nil
}

// GetGradeList 获取年级列表
func (ctr AnswerController) GetGradeList(ctx *gin.Context) (gin.H, error) {
	gradeList := activity.DefaultGDdbService.GetGradeList()
	return gin.H{
		"record": gradeList,
	}, nil
}

func (ctr AnswerController) CreateSchool(ctx *gin.Context) (gin.H, error) {
	form := &GDDbCreateSchoolForm{}
	if err := apiutil.BindForm(ctx, form); err != nil {
		return nil, err
	}
	id, err := activity.DefaultGDdbService.CreateSchool(form.SchoolName, form.CityId, form.GradeType)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"schoolId": id,
	}, nil
}

// GetSchoolList 获取学校列表
func (ctr AnswerController) GetSchoolList(ctx *gin.Context) (gin.H, error) {
	form := &GDDbSelectSchoolForm{}
	if err := apiutil.BindForm(ctx, form); err != nil {
		return nil, err
	}
	schoolList := activity.DefaultGDdbService.GetSchoolList(form.SchoolName, form.CityId, form.GradeId)
	return gin.H{
		"record": schoolList,
	}, nil
}

// GetAchievement 获取我的成就
func (ctr AnswerController) GetAchievement(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	res := activity.DefaultGDdbService.GetAchievement(user.ID)
	return gin.H{
		"titleUrl":       res.TitleUrl,
		"certificateUrl": res.CertificateUrl,
	}, nil
}

// GetIntegral 获取积分
func (ctr AnswerController) GetIntegral(ctx *gin.Context) (gin.H, error) {
	return nil, nil
}

// CloseLateTips 关闭提示
func (ctr AnswerController) CloseLateTips(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	if err := activity.DefaultGDdbService.CloseLateTips(user.ID); err != nil {
		return nil, err
	}
	return nil, nil
}
