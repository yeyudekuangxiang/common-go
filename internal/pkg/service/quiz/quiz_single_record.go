package quiz

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/util/timeutils"
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
	t := model.NewTime()
	record := entity.QuizSingleRecord{
		OpenId:     param.OpenId,
		QuestionId: param.QuestionId,
		Correct:    param.Correct,
		AnswerTime: t,
		AnswerDate: t.Date(),
	}
	return &record, app.DB.Create(&record).Error
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
