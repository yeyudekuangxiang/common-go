package request

type ApiUserId struct {
	ID int64 `form:"id" binding:"gt=0" alias:"用户id"`
}
