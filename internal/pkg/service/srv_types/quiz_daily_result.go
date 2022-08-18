package srv_types

import "mio/internal/pkg/model/entity"

type QuizDailyResult struct {
	entity.QuizDailyResult
	Text []string `json:"text"`
}
