package quiz

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"time"
)

const OneDayAnswerNum = 4
const questionToPointRatio = 25

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
	if !util.DefaultLock.Lock("QuizAnswerQuestion"+openid, time.Second*5) {
		return nil, errno.ErrLimit
	}
	defer util.DefaultLock.UnLock("QuizAnswerQuestion" + openid)

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
	if !util.DefaultLock.Lock(fmt.Sprintf("QUIZ_Ssubmit%s", openId), time.Second*10) {
		return errno.ErrLimit
	}
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

	_, err := service.NewPointService(context.NewMioContext()).IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:      openId,
		Type:        entity.POINT_QUIZ,
		BizId:       util.UUID(),
		ChangePoint: int64(correctNum * questionToPointRatio),
	})
	return err
}
func (srv QuizService) DailyResult(openId string) (*srv_types.QuizDailyResult, error) {
	result, err := DefaultQuizDailyResultService.FindTodayResult(openId)
	if err != nil {
		return nil, err
	}
	qdResult := srv_types.QuizDailyResult{
		QuizDailyResult: *result,
	}
	switch result.CorrectNum {
	case 0:
		qdResult.Text = []string{"妥妥的一只零碳萌新~", "还需努力哟！"}
	case 1:
		qdResult.Text = []string{"零碳见习生你好呀~", "要再接再厉嗷~"}
	case 2:
		qdResult.Text = []string{"你已达到零碳高中生水平！", "加油吧少年~"}
	case 3:
		qdResult.Text = []string{"你已成功晋级零碳大学生！", "距离巅峰只有一步之遥~"}
	case 4:
		qdResult.Text = []string{"你已达到零碳博士水准~", "再创辉煌吧！"}
	}
	return &qdResult, nil
}
func (srv QuizService) GetSummary(openId string) (*entity.QuizSummary, error) {
	return DefaultQuizSummaryService.FindOrCreateSummary(openId)
}
