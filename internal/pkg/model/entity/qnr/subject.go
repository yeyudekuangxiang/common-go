package qnr

type Subject struct {
	ID         int64
	SubjectId  int64
	QnrId      int64
	CategoryId int64
	Title      string
	Type       int8
	IsHide     int8
	Remind     string
	Sort       int8
}

func (Subject) TableName() string {
	return "qnr_subject"
}
