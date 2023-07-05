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
	result := make([]QuizQuestionAPI, 0, len(list))
	for _, item := range list {
		result = append(result, QuizQuestionAPI{
			ID:                  strconv.FormatInt(item.ID, 10),
			QuestionStatement:   item.QuestionStatement,
			Choices:             randomOptions(item.Choices),
			AnswerStatement:     item.AnswerStatement,
			DetailedDescription: item.DetailedDescription,
			Type:                item.Type,
			QuestionID:          strconv.FormatInt(item.ID, 10),
		})
	}
	return gin.H{
		"list": result,
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
