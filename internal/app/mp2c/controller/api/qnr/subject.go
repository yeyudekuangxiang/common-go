package qnr

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/context"
	qnrService "mio/internal/pkg/service/qnr"
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
	form := api_types.GetQnrSubjectCreateDTO{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)
	err := answerServer.Add(srv_types.AddQnrAnswerDTO{
		OpenId: user.OpenId,
		UserId: user.ID,
		Answer: form.Answer})
	if err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (SubjectController) GetList(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	subjectServer := qnrService.NewSubjectService(context.NewMioContext(context.WithContext(ctx)))
	ret, err := subjectServer.GetList(user.OpenId)
	return ret, err
}
