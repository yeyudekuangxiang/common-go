package qnr

type Answer struct {
	ID        int64
	QnrId     int64
	SubjectId int64
	UserId    int64
	Answer    string
}

func (Answer) TableName() string {
	return "qnr_answer"
}
