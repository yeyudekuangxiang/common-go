package api

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"math/rand"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/service/quiz"
	"mio/internal/pkg/util/apiutil"
	"mio/pkg/errno"
	"sort"
	"strconv"
	"strings"
)

var DefaultQuizController = QuizController{}

type QuizController struct {
}

func (QuizController) GetDailyQuestions(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	form := GetDailyQuestionsForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	list, err := quiz.DefaultQuizService.DailyQuestions(user.OpenId, form.ActivityChannel)
	result := make([]QuizQuestionAPI, 0, len(list))
	for _, item := range list {
		var choice []string
		err := json.Unmarshal([]byte(item.Choices), &choice)
		if err != nil {
			return nil, err
		}
		result = append(result, QuizQuestionAPI{
			ID:                  strconv.FormatInt(item.ID, 10),
			QuestionStatement:   item.QuestionStatement,
			Choices:             randomOptions(choice),
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
	isBlack := app.Redis.SIsMember(context.Background(), config.RedisKey.QuizBlackList, user.OpenId)
	status := int64(0)
	desc := ""
	if isBlack.Val() {
		status = 1
		desc = "已加入答题黑名单"
	}
	return gin.H{
		"availability": availability,
		"status":       status, // 0正常 1 答题黑名单
		"desc":         desc,
	}, err
}

func (QuizController) AddAvailability(ctx *gin.Context) (gin.H, error) {
	form := AddAvailabilityReq{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	if form.Openid == "" || form.Key != "liumei_init" {
		return gin.H{
			"status": "Openid未空或者Key有误",
		}, nil
	}
	openIds := strings.Split(form.Openid, ",")
	for _, openid := range openIds {
		if form.Type == "add" {
			ret := app.Redis.SAdd(context.Background(), config.RedisKey.QuizBlackList, openid)
			if ret.Err() != nil {
				return gin.H{
					"status": "redis入库失败",
				}, nil
			}
		}
		if form.Type == "sub" {
			ret := app.Redis.SRem(context.Background(), config.RedisKey.QuizBlackList, openid)
			if ret.Err() != nil {
				return gin.H{
					"status": "redis入库失败",
				}, nil
			}
		}

	}
	members := app.Redis.SMembers(context.Background(), config.RedisKey.QuizBlackList)

	return gin.H{
		"status": members.Val(),
	}, nil
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
	isBlack := app.Redis.SIsMember(context.Background(), config.RedisKey.QuizBlackList, user.OpenId)
	if isBlack.Val() {
		return nil, errno.ErrCommon.WithMessage("您的账号异常，无法参与答题")
	}
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
