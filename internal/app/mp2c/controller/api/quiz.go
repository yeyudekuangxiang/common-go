package api

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"mio/internal/pkg/service/quiz"
	"mio/internal/pkg/util/apiutil"
	"sort"
	"strconv"
)

var DefaultQuizController = QuizController{}

type QuizController struct {
}

func (QuizController) GetDailyQuestions(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)

	list, err := quiz.DefaultQuizService.DailyQuestions(user.OpenId)

	for i, item := range list {
		list[i].Choices = randomOptions(item.Choices)
		list[i].QuestionId = strconv.FormatInt(list[i].ID, 10)
	}
	return gin.H{
		"list": list,
	}, err
}
func randomOptions(options []string) []string {
	sort.Slice(options, func(i, j int) bool {
		return rand.Intn(2) == 0
	})
	return options
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
	point, err := quiz.DefaultQuizService.Submit(user.OpenId, user.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"point": point,
	}, nil
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
