package quiz

import (
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultQuizQuestionService = QuizQuestionService{}

type QuizQuestionService struct {
}

func (srv QuizQuestionService) GetDailyQuestions(num int) ([]entity.QuizQuestion, error) {
	list := make([]entity.QuizQuestion, 0)
	err := app.DB.Order("random() asc").Limit(num).Find(&list).Error
	if err != nil {
		panic(err)
	}
	return list, nil
}
