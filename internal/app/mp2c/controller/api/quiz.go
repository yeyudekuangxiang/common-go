package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/service/quiz"
	"mio/internal/pkg/util/apiutil"
)

var DefaultQuizController = QuizController{}

type QuizController struct {
}

func (QuizController) GetDailyQuestions(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)

	list, err := quiz.DefaultQuizService.DailyQuestions(user.OpenId)
	return gin.H{
		"list": list,
	}, err
}
func (QuizController) Availability(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	availability, err := quiz.DefaultQuizService.Availability(user.OpenId)

	return gin.H{
		"availability": availability,
	}, err
}

func (QuizController) AnswerQuestion(ctx *gin.Context) (gin.H, error) {
	form := AnswerQuizQuestionForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)
	result, err := quiz.DefaultQuizService.AnswerQuestion(user.OpenId, form.QuestionId, form.Choice)
	return gin.H{
		"answerResult": result,
	}, err
}
func (QuizController) Submit(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	return nil, quiz.DefaultQuizService.Submit(user.OpenId)
}

func (QuizController) DailyResult(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	result, err := quiz.DefaultQuizService.DailyResult(user.OpenId)
	return gin.H{
		"dailyResult": result,
	}, err
}
func (QuizController) GetSummary(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	summary, err := quiz.DefaultQuizService.GetSummary(user.OpenId)
	return gin.H{
		"summary": summary,
	}, err
}
