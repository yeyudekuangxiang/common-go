package entity

import (
	"database/sql/driver"
	"github.com/pkg/errors"
	"strings"
)

type Choices []string

func (d Choices) Value() (driver.Value, error) {
	build := strings.Builder{}
	build.WriteString("{")
	if d != nil {
		build.WriteString(strings.Join(d, ","))
	}
	build.WriteString("}")
	return build.String(), nil
}
func (d *Choices) Scan(value interface{}) error {
	t, ok := value.(string)
	if !ok {
		return errors.New("Choices type error")
	}
	t = t[1 : len(t)-1]
	*d = strings.Split(t, ",")
	return nil
}

type QuizQuestion struct {
	ID                  int64   `json:"id"`
	QuestionId          string  `json:"questionId"`
	QuestionStatement   string  `json:"questionStatement"`
	Choices             Choices `json:"choices"`
	AnswerStatement     string  `json:"answerStatement"`
	DetailedDescription string  `json:"detailedDescription"`
}
