package quiz

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/util/timeutils"
	"mio/pkg/errno"
)

var DefaultQuizDailyResultService = QuizDailyResultService{}

type QuizDailyResultService struct {
}

// IsAnsweredToday 今天是否已经完成答题 true已完成答题 false 未完成答题
func (srv QuizDailyResultService) IsAnsweredToday(openId string, day timeutils.Date) (bool, error) {
	result, err := srv.FindTodayResult(openId, day)
	return result.ID != 0, err
}
func (srv QuizDailyResultService) FindTodayResult(openId string, day timeutils.Date) (*entity.QuizDailyResult, error) {
	result := entity.QuizDailyResult{}
	err := app.DB.Where("openid = ? and answer_time >= ? and answer_time< ?", openId, day.FullString(), day.AddDay(1).FullString()).Take(&result).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return &result, nil
}
func (srv QuizDailyResultService) CompleteTodayQuiz(openid string, t timeutils.Time) (*entity.QuizDailyResult, error) {
	todayResult, err := srv.FindTodayResult(openid, t.Date())
	if err != nil {
		return nil, err
	}
	if todayResult.ID != 0 {
		return nil, errno.ErrCommon.WithMessage("请勿重复提交答题")
	}
	todaySummary := DefaultQuizSingleRecordService.GetTodaySummary(openid, t.Date())
	if todaySummary.AnsweredNum < OneDayAnswerNum {
		return nil, errno.ErrCommon.WithMessage("答题未完成")
	}

	todayResult = &entity.QuizDailyResult{
		OpenId:       openid,
		AnswerDate:   t.Date().Time,
		AnswerTime:   t.Time,
		CorrectNum:   int(todaySummary.CorrectNum),
		IncorrectNum: int(todaySummary.AnsweredNum - todaySummary.CorrectNum),
	}
	return todayResult, app.DB.Create(todayResult).Error
}
