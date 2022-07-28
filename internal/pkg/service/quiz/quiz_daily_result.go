package quiz

import (
	"errors"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/util/timeutils"
)

var DefaultQuizDailyResultService = QuizDailyResultService{}

type QuizDailyResultService struct {
}

// IsAnsweredToday 今天是否已经完成答题 true已完成答题 false 未完成答题
func (srv QuizDailyResultService) IsAnsweredToday(openId string) (bool, error) {
	result, err := srv.FindTodayResult(openId)
	return result.ID != 0, err
}
func (srv QuizDailyResultService) FindTodayResult(openId string) (*entity.QuizDailyResult, error) {
	result := entity.QuizDailyResult{}
	err := app.DB.Where("openid = ? and answer_date = ?", openId, timeutils.NowDate().FullString()).Take(&result).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return &result, nil
}
func (srv QuizDailyResultService) CompleteTodayQuiz(openid string) (*entity.QuizDailyResult, error) {
	todayResult, err := srv.FindTodayResult(openid)
	if err != nil {
		return nil, err
	}
	if todayResult.ID != 0 {
		return nil, errors.New("请勿重复提交答题")
	}
	todaySummary := DefaultQuizSingleRecordService.GetTodaySummary(openid)
	if todaySummary.AnsweredNum < OneDayAnswerNum {
		return nil, errors.New("答题未完成")
	}

	t := model.NewTime()
	todayResult = &entity.QuizDailyResult{
		OpenId:       openid,
		AnswerDate:   t.Date(),
		AnswerTime:   t,
		CorrectNum:   int(todaySummary.CorrectNum),
		IncorrectNum: int(todaySummary.AnsweredNum - todaySummary.CorrectNum),
	}
	return todayResult, app.DB.Create(todayResult).Error
}
