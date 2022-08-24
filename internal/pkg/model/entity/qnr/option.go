package qnr

type Option struct {
	ID             int64
	Title          string
	Sort           int8
	SubjectId      int64
	Remind         string
	JumpSubject    int64
	RelatedSubject string
}

func (Option) TableName() string {
	return "qnr_option"
}
