package quiz

import (
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/timetool"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/pkg/errno"
	"time"
)

var DefaultQuizSingleRecordService = QuizSingleRecordService{}

type QuizSingleRecordService struct {
}

func (srv QuizSingleRecordService) ClearTodayRecord(openid string) {
	err := app.DB.Where("openid = ? and answer_time >= ? and answer_time < ?", openid, timetool.NowDate().FullString(), timetool.NowDate().AddDay(1).FullString()).Delete(&entity.QuizSingleRecord{}).Error
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
		AnswerDate: timetool.NowDate().Time,
	}
	return &record, app.DB.Create(&record).Error
}
func (srv QuizSingleRecordService) IsAnswered(openId string, questionId string, day timetool.Date) error {
	record := entity.QuizSingleRecord{}
	err := app.DB.Where("openid = ? and question_id = ? and answer_time >= ? and answer_time < ?", openId, questionId, day.FullString(), day.AddDay(1).FullString()).Take(&record).Error

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
	err := app.DB.Model(entity.QuizSingleRecord{}).Where("openid = ? and answer_time >= ? and answer_time < ?", openId, timetool.NowDate().FullString(), timetool.NowDate().AddDay(1).FullString()).Count(&count).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return count
}
func (srv QuizSingleRecordService) GetTodaySummary(openId string, day timetool.Date) DaySummary {
	summary := DaySummary{}
	err := app.DB.Model(entity.QuizSingleRecord{}).Where("openid = ? and answer_time >= ? and answer_time < ?", openId, day.FullString(), day.AddDay(1).FullString()).Count(&summary.AnsweredNum).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	err = app.DB.Model(entity.QuizSingleRecord{}).Where("openid = ? and answer_time >= ? and answer_time < ? and correct = true", openId, day.FullString(), day.AddDay(1).FullString()).Count(&summary.CorrectNum).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return summary
}
