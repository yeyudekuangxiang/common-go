package question

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/context"
	qnrService "mio/internal/pkg/service/question"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util/apiutil"
)

var DefaultSubjectController = SubjectController{}

type SubjectController struct {
}

type Ans struct {
	Id     int64  `json:"id"`
	Answer string `json:"answer"`
}

func (SubjectController) Create(ctx *gin.Context) (gin.H, error) {
	answerServer := qnrService.NewAnswerService(context.NewMioContext(context.WithContext(ctx)))
	form := api_types.GetQuestionSubjectCreateDTO{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)
	if user.PhoneNumber == "" {
		//	return gin.H{}, errno.ErrCommon.WithMessage("请您先绑定手机号")
	}
	err := answerServer.Add(srv_types.AddQuestionAnswerDTO{
		OpenId:     user.OpenId,
		UserId:     user.ID,
		Answer:     form.Answer,
		QuestionId: 1})
	if err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (SubjectController) GetList(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	if user.PhoneNumber == "" {
		//return gin.H{}, errno.ErrCommon.WithMessage("请您先绑定手机号")
	}
	subjectServer := qnrService.NewSubjectService(context.NewMioContext(context.WithContext(ctx)))
	ret, err := subjectServer.GetList(user.OpenId, 1)
	return ret, err
}

func (receiver SubjectController) GetUserQuestion(ctx *gin.Context) (gin.H, error) {
	subjectServer := qnrService.NewSubjectService(context.NewMioContext(context.WithContext(ctx)))
	//获取问卷碳量
	ret := subjectServer.GetUserQuestion(srv_types.GetQuestionUserDTO{UserId: 1570653152666181632, QuestionId: 1})

	return gin.H{"userCarbonInfo": ret}, nil
}
