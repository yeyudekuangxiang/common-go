package types

type User struct {
	//用户id
	Id int `json:"id"`
	//用户昵称
	Name string `json:"name"`
	//用户账号
	Username string `json:"username"`
	//用户头像
	AvatarUrl string `json:"avatar_url"`
	//用户邮箱
	Email string `json:"email"`
}
