package activity

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/service/activity"
	"mio/internal/pkg/util/apiutil"
)

var DefaultAnswerController = AnswerController{}

type AnswerController struct {
}

// HomePage Answer the questions
func (ctr AnswerController) HomePage(ctx *gin.Context) (gin.H, error) {
	//当前登陆用户信息
	user := apiutil.GetAuthUser(ctx)
	form := &GDDbActivityHomePageForm{}
	if err := apiutil.BindForm(ctx, form); err != nil {
		return nil, err
	}

	res, err := activity.DefaultGDdbService.HomePage(user.ID, form.UserId)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"record": res,
	}, err
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
		return nil, err
	}
	// 用户用户身份更新状态及学校排名
	err = activity.DefaultGDdbService.CheckActivityStatus(user.ID, form.SchoolId)
	if err != nil {
		//用户已更新为团长，生成称号
		return nil, err
	}
	// 受邀者 or 邀请者 答题完成后赠书+1
	err = activity.DefaultGDdbService.IncrRank(form.SchoolId)
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
	schoolList := activity.DefaultGDdbService.CreateSchool(form.SchoolName, form.CityId, form.GradeType)
	return gin.H{
		"record": schoolList,
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
func (ctr AnswerController) GetIntegral(ctx *gin.Context) error {
	return nil
}
