package quiz

type CreateSingleRecordParam struct {
	OpenId     string
	QuestionId string
	Correct    bool
}
type AnswerQuestionResult struct {
	IsMatched           bool   `json:"isMatched"`
	DetailedDescription string `json:"detailedDescription"`
	CurrentIndex        int64  `json:"currentIndex"`
}
type UpdateSummaryParam struct {
	OpenId           string
	TodayCorrectNum  int
	TodayAnsweredNum int
}
type DaySummary struct {
	AnsweredNum int64
	CorrectNum  int64
}
