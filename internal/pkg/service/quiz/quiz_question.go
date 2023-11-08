package quiz

import (
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultQuizQuestionService = QuizQuestionService{}

type QuizQuestionService struct {
}

func (srv QuizQuestionService) GetDailyQuestions(num int) ([]entity.QuizQuestionV2, error) {
	list := make([]entity.QuizQuestionV2, 0)
	err := app.DB.Order("random()").Where("channel = ?", 0).Limit(num).Find(&list).Error
	if err != nil {
		panic(err)
	}
	return list, nil
}

func (srv QuizQuestionService) GetDailyQuestionsById(num int, id int64) ([]entity.QuizQuestionV2, error) {
	list := make([]entity.QuizQuestionV2, 0)
	err := app.DB.Order("random()").Where("id = ?", id).Limit(num).Find(&list).Error
	if err != nil {
		panic(err)
	}
	return list, nil
}
