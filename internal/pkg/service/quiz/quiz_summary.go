package quiz

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
)

var DefaultQuizSummaryService = QuizSummaryService{}

type QuizSummaryService struct {
}

func (srv QuizSummaryService) UpdateTodaySummary(param UpdateSummaryParam) error {
	summary, err := srv.FindOrCreateSummary(param.OpenId)
	if err != nil {
		return err
	}

	summary.TotalAnsweredNum += param.TodayAnsweredNum
	summary.TotalCorrectNum += param.TodayCorrectNum
	return app.DB.Save(&summary).Error
}
func (srv QuizSummaryService) FindOrCreateSummary(openId string) (*entity.QuizSummary, error) {
	summary := entity.QuizSummary{}
	err := app.DB.Where("openid = ?", openId).Take(&summary).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	if summary.ID != 0 {
		return &summary, nil
	}

	summary.OpenId = openId
	t := model.NewTime()
	summary.LastUpdateDate = t.Date()
	err = app.DB.Create(&summary).Error
	if err != nil {
		return nil, err
	}

	return &summary, nil
}
