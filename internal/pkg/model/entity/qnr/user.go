package qnr

type User struct {
	UserId      int64
	ThirdId     string
	InvitedById string
	Phone       string
	AddressId   string
	channel     string
	Ip          string
	City        string
}

func (User) TableName() string {
	return "qnr_user"
}
