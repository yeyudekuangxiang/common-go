package entity

type QuizQuestionV2 struct {
	ID                  int64  `json:"id"`
	QuestionId          string `json:"questionId"`
	QuestionStatement   string `json:"questionStatement"`
	Choices             string `json:"choices"`
	AnswerStatement     string `json:"answerStatement"`
	DetailedDescription string `json:"detailedDescription"`
	Type                int    `json:"type"`
	Channel             string `json:"channel"`
}

func (m QuizQuestionV2) TableName() string {
	return "quiz_question_v2"
}
