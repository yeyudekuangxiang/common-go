package quiz

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/util/timeutils"
	"mio/pkg/errno"
	"time"
)

var DefaultQuizSingleRecordService = QuizSingleRecordService{}

type QuizSingleRecordService struct {
}

func (srv QuizSingleRecordService) ClearTodayRecord(openid string) {
	err := app.DB.Where("openid = ? and answer_time >= ? and answer_time < ?", openid, timeutils.NowDate().FullString(), timeutils.NowDate().AddDay(1).FullString()).Delete(&entity.QuizSingleRecord{}).Error
	if err != nil {
		panic(err)
	}
	return
}
func (srv QuizSingleRecordService) CreateSingleRecord(param CreateSingleRecordParam) (*entity.QuizSingleRecord, error) {

	record := entity.QuizSingleRecord{
		OpenId:     param.OpenId,
		QuestionId: param.QuestionId,
		Correct:    param.Correct,
		AnswerTime: time.Now(),
		AnswerDate: timeutils.NowDate().Time,
	}
	return &record, app.DB.Create(&record).Error
}
func (srv QuizSingleRecordService) IsAnswered(openId string, questionId string) error {
	record := entity.QuizSingleRecord{}
	err := app.DB.Where("openid = ? and question_id = ? and answer_time >= ? and answer_time < ?", openId, questionId, timeutils.NowDate().FullString(), timeutils.NowDate().AddDay(1).FullString()).Take(&record).Error

	if err == nil {
		return errno.ErrLimit.WithMessage("重复提交")
	}

	if err == gorm.ErrRecordNotFound {
		return nil
	}

	return err
}

// GetTodayAnswerNum 获取今天已答题数量
func (srv QuizSingleRecordService) GetTodayAnswerNum(openId string) int64 {
	var count int64
	err := app.DB.Model(entity.QuizSingleRecord{}).Where("openid = ? and answer_time >= ? and answer_time < ?", openId, timeutils.NowDate().FullString(), timeutils.NowDate().AddDay(1).FullString()).Count(&count).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return count
}
func (srv QuizSingleRecordService) GetTodaySummary(openId string) DaySummary {
	summary := DaySummary{}
	err := app.DB.Model(entity.QuizSingleRecord{}).Where("openid = ? and answer_time >= ? and answer_time < ?", openId, timeutils.NowDate().FullString(), timeutils.NowDate().AddDay(1).FullString()).Count(&summary.AnsweredNum).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	err = app.DB.Model(entity.QuizSingleRecord{}).Where("openid = ? and answer_time >= ? and answer_time < ? and correct = true", openId, timeutils.NowDate().FullString(), timeutils.NowDate().AddDay(1).FullString()).Count(&summary.CorrectNum).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return summary
}
