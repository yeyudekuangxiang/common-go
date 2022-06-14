package quiz

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util"
	"time"
)

const OneDayAnswerNum = 4
const questionToPointRatio = 50

var DefaultQuizService = QuizService{}

type QuizService struct {
}

func (srv QuizService) DailyQuestions(openid string) ([]entity.QuizQuestion, error) {
	DefaultQuizSingleRecordService.ClearTodayRecord(openid)
	return DefaultQuizQuestionService.GetDailyQuestions(OneDayAnswerNum)
}

// Availability 是否可以答题 true可以 false不可以
func (srv QuizService) Availability(openid string) (bool, error) {
	isAnsweredToday, err := DefaultQuizDailyResultService.IsAnsweredToday(openid)
	if err != nil {
		return false, err
	}
	return !isAnsweredToday, nil
}
func (srv QuizService) AnswerQuestion(openid, questionId, answer string) (*AnswerQuestionResult, error) {
	todayAnsweredNum := DefaultQuizSingleRecordService.GetTodayAnswerNum(openid)
	if todayAnsweredNum >= OneDayAnswerNum {
		return nil, errors.New("答题数量超出限制")
	}

	question := entity.QuizQuestion{}

	err := app.DB.Where("question_id = ?", questionId).Take(&question).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	if question.ID == 0 {
		return nil, errors.New("题目不存在")
	}

	isRight := question.AnswerStatement == answer
	_, err = DefaultQuizSingleRecordService.CreateSingleRecord(CreateSingleRecordParam{
		OpenId:     openid,
		QuestionId: questionId,
		Correct:    isRight,
	})
	if err != nil {
		return nil, err
	}
	return &AnswerQuestionResult{
		IsMatched:           isRight,
		DetailedDescription: question.DetailedDescription,
		CurrentIndex:        DefaultQuizSingleRecordService.GetTodayAnswerNum(openid),
	}, nil
}
func (srv QuizService) Submit(openId string) error {
	util.DefaultLock.Lock(fmt.Sprintf("QUIZ_Ssubmit%s", openId), time.Second*5)
	defer util.DefaultLock.UnLock(fmt.Sprintf("QUIZ_Ssubmit%s", openId))

	todayResult, err := DefaultQuizDailyResultService.CompleteTodayQuiz(openId)
	if err != nil {
		return err
	}
	err = DefaultQuizSummaryService.UpdateTodaySummary(UpdateSummaryParam{
		OpenId:           openId,
		TodayCorrectNum:  todayResult.CorrectNum,
		TodayAnsweredNum: todayResult.IncorrectNum + todayResult.CorrectNum,
	})
	if err != nil {
		return err
	}

	return srv.SendAnswerPoint(openId, todayResult.CorrectNum)
}
func (srv QuizService) SendAnswerPoint(openId string, correctNum int) error {
	if correctNum > OneDayAnswerNum {
		correctNum = OneDayAnswerNum
	}

	_, err := service.DefaultPointTransactionService.Create(service.CreatePointTransactionParam{
		OpenId: openId,
		Type:   entity.POINT_QUIZ,
		Value:  correctNum * questionToPointRatio,
	})
	return err
}
func (srv QuizService) DailyResult(openId string) (*entity.QuizDailyResult, error) {
	return DefaultQuizDailyResultService.FindTodayResult(openId)
}
func (srv QuizService) GetSummary(openId string) (*entity.QuizSummary, error) {
	return DefaultQuizSummaryService.FindOrCreateSummary(openId)
}
