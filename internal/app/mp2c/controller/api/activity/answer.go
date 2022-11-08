package activity

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/app"
	activity3 "mio/internal/pkg/model/entity/activity"
	activity2 "mio/internal/pkg/repository/activity"
	"mio/internal/pkg/service/activity"
	"mio/internal/pkg/service/oss"
	"mio/internal/pkg/util/apiutil"
	"mio/pkg/errno"
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
	//初始化response
	homeUser := activity2.GDDbHomePageUserInfo{
		UserInfo:    activity2.GDDbUserInfo{},
		InviteInfo:  activity2.GDDbUserInfo{},
		InvitedInfo: make([]activity2.GDDbUserInfo, 0),
	}
	homeSchool := make([]activity3.GDDbSchoolRank, 0)
	var err error
	var isNewUser bool
	//获取form
	form := &GDDbActivityHomePageForm{}
	if err = apiutil.BindForm(ctx, form); err != nil {
		return nil, err
	}
	if user.ID != 0 {
		if activity.DefaultGDdbService.IsNewUser(user) {
			isNewUser = true
		}
		homeUser, err = activity.DefaultGDdbService.HomePageUser(user.ID, form.InviteId, isNewUser)
		if err != nil {
			return nil, err
		}
	}
	homeSchool, _ = activity.DefaultGDdbService.HomePageSchool()
	return gin.H{
		"record": activity.GDDbHomePageResponse{
			User:      homeUser,
			School:    homeSchool,
			IsNewUser: isNewUser,
		},
	}, nil
}

// GetUserSchool 获取用户及学校信息
func (ctr AnswerController) GetUserSchool(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	if user.ID == 0 {
		return nil, errno.ErrCommon.WithMessage("未登录，无法访问。")
	}
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
	if user.ID == 0 {
		return nil, errno.ErrCommon.WithMessage("未登录，无法访问。")
	}
	if !activity.DefaultGDdbService.IsNewUser(user) {
		return nil, errno.ErrCommon.WithMessage("本次活动仅限绿喵新用户参加哦～")
	}
	t, _ := strconv.Atoi(ctx.PostForm("type"))
	file, err := ctx.FormFile("file")
	if err != nil {
		return nil, err
	}
	if file.Size > 5*1024*1024 {
		return nil, errno.ErrCommon.WithMessage("上传失败，请调整文件大小！")
	}
	open, err := file.Open()
	if err != nil {
		return nil, err
	}
	object, err := oss.DefaultOssService.PutObject(fmt.Sprintf("/activity/gd/%d_%d", user.ID, time.Now().Unix()), open)
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
	if user.ID == 0 {
		return nil, errno.ErrCommon.WithMessage("未登录，无法访问。")
	}
	if !activity.DefaultGDdbService.IsNewUser(user) {
		return nil, errno.ErrCommon.WithMessage("本次活动仅限绿喵新用户参加哦～")
	}
	return nil, nil
}

//EndQuestion 答题完成
func (ctr AnswerController) EndQuestion(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	if user.ID == 0 {
		return nil, errno.ErrCommon.WithMessage("未登录，无法访问。")
	}
	if !activity.DefaultGDdbService.IsNewUser(user) {
		return nil, errno.ErrCommon.WithMessage("本次活动仅限绿喵新用户参加哦～")
	}
	form := &GDDbActivitySchoolForm{}
	if err := apiutil.BindForm(ctx, form); err != nil {
		return nil, err
	}

	// 保存学校信息，更新答题状态
	app.Logger.Info("答题完成，保存学校信息start")
	err := activity.DefaultGDdbService.SaveSchoolInfo(form.UserName, form.SchoolId, form.GradeId, user.ID, form.ClassNumber)
	app.Logger.Info("答题完成，保存学校信息end")

	if err != nil {
		return nil, errno.ErrCommon.WithMessage("学校信息保存失败")
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
	user := apiutil.GetAuthUser(ctx)
	if user.ID == 0 {
		return nil, errno.ErrCommon.WithMessage("未登录，无法访问。")
	}
	if !activity.DefaultGDdbService.IsNewUser(user) {
		return nil, errno.ErrCommon.WithMessage("本次活动仅限绿喵新用户参加哦～")
	}
	cityList := activity.DefaultGDdbService.GetCityList()
	return gin.H{
		"record": cityList,
	}, nil
}

// GetGradeList 获取年级列表
func (ctr AnswerController) GetGradeList(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	if user.ID == 0 {
		return nil, errno.ErrCommon.WithMessage("未登录，无法访问。")
	}
	if !activity.DefaultGDdbService.IsNewUser(user) {
		return nil, errno.ErrCommon.WithMessage("本次活动仅限绿喵新用户参加哦～")
	}
	gradeList := activity.DefaultGDdbService.GetGradeList()
	return gin.H{
		"record": gradeList,
	}, nil
}

func (ctr AnswerController) CreateSchool(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	if user.ID == 0 {
		return nil, errno.ErrCommon.WithMessage("未登录，无法访问。")
	}
	if !activity.DefaultGDdbService.IsNewUser(user) {
		return nil, errno.ErrCommon.WithMessage("本次活动仅限绿喵新用户参加哦～")
	}
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
	user := apiutil.GetAuthUser(ctx)
	if user.ID == 0 {
		return nil, errno.ErrCommon.WithMessage("未登录，无法访问。")
	}
	if !activity.DefaultGDdbService.IsNewUser(user) {
		return nil, errno.ErrCommon.WithMessage("本次活动仅限绿喵新用户参加哦～")
	}
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
	if user.ID == 0 {
		return nil, errno.ErrCommon.WithMessage("未登录，无法访问。")
	}
	if !activity.DefaultGDdbService.IsNewUser(user) {
		return nil, errno.ErrCommon.WithMessage("本次活动仅限绿喵新用户参加哦～")
	}
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
	if user.ID == 0 {
		return nil, errno.ErrCommon.WithMessage("未登录，无法访问。")
	}
	if !activity.DefaultGDdbService.IsNewUser(user) {
		return nil, errno.ErrCommon.WithMessage("本次活动仅限绿喵新用户参加哦～")
	}
	if err := activity.DefaultGDdbService.CloseLateTips(user.ID); err != nil {
		return nil, err
	}
	return nil, nil
}
