package srv_types

import "mio/internal/pkg/model/entity"

type QuizDailyResult struct {
	entity.QuizDailyResult
	Point int      `json:"point"`
	Text  []string `json:"text"`
}
