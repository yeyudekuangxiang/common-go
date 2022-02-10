package request

type ApiUserId struct {
	ID int `form:"id" binding:"gt=0" alias:"用户id"`
}
