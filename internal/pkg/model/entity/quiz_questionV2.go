package entity

type QuizQuestionV2 struct {
	ID                  int64   `json:"id"`
	QuestionId          string  `json:"questionId"`
	QuestionStatement   string  `json:"questionStatement"`
	Choices             Choices `json:"choices"`
	AnswerStatement     string  `json:"answerStatement"`
	DetailedDescription string  `json:"detailedDescription"`
	Type                int     `json:"type"`
}

func (m QuizQuestionV2) TableName() string {
	return "quiz_questions_v2"
}
