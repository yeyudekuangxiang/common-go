package quizmsg

type QuizMessage struct {
	OpenId           string
	QuizTime         int64
	TodayCorrectNum  int
	TodayAnsweredNum int
}

type QuizSendMessage struct {
	Uid    int64
	OpenId string
	BizId  string
}
