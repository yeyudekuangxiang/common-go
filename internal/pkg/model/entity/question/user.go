package question

type User struct {
	UserId      int64
	ThirdId     string
	InvitedById string
	Phone       string
	Channel     string
	Ip          string
	City        string
}

func (User) TableName() string {
	return "question_user"
}
